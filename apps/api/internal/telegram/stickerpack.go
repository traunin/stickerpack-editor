package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/config"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/retrier"
)

var (
	httpClient   = &http.Client{Timeout: 15 * time.Second}
	fetchRetires = 3
	ErrorPackNotFound = errors.New("pack was deleted or does not exist")
)

type InputSticker struct {
	Sticker   []byte
	Format    string
	EmojiList []string
	Keywords  []string
}

type PackPreview struct {
	Title       string `json:"title"`
	Name        string `json:"name"`
	ThumbnailID string `json:"thumbnail_id"`
}

type StickerPack struct {
	userID      int64
	name        string
	title       string
	stickers    []InputSticker
	isPublic    bool
	thumbnailID string

	nameSet      bool
	validNameSet bool
}

type attachSticker struct {
	Sticker   string   `json:"sticker"`
	Format    string   `json:"format"`
	EmojiList []string `json:"emoji_list"`
	Keywords  []string `json:"keywords"`
}

type GetStickerSetResponse struct {
	Ok          bool       `json:"ok"`
	Result      StickerSet `json:"result,omitempty"`
	Description string     `json:"description,omitempty"`
}

func (sp *StickerPack) UserID() int64 {
	return sp.userID
}

func (sp *StickerPack) Name() string {
	return sp.name
}

func (sp *StickerPack) Title() string {
	return sp.title
}

func (sp *StickerPack) IsPublic() bool {
	return sp.isPublic
}

func (sp *StickerPack) ThumbnailID() string {
	return sp.thumbnailID
}

type Option func(*StickerPack)

func WithValidName(validName string) Option {
	return func(sp *StickerPack) {
		sp.validNameSet = true
		sp.name = validName
	}
}

func WithName(name string) Option {
	return func(sp *StickerPack) {
		sp.nameSet = true
		validName := ValidPackName(name)
		sp.name = validName
	}
}

func WithTitle(title string) Option {
	return func(sp *StickerPack) {
		sp.title = title
	}
}

func WithStickers(stickers []InputSticker) Option {
	return func(sp *StickerPack) {
		sp.stickers = stickers
	}
}

func WithPublic(public bool) Option {
	return func(sp *StickerPack) {
		sp.isPublic = public
	}
}

func ValidPackName(name string) string {
	botName := config.Load().BotName()
	return fmt.Sprintf("%s_by_%s", name, botName)
}

func NewStickerPack(userID int64, opts ...Option) (*StickerPack, error) {
	sp := &StickerPack{userID: userID}
	for _, opt := range opts {
		opt(sp)
	}

	if sp.nameSet && sp.validNameSet {
		return nil, fmt.Errorf("cannot use both WithName and WithValidName")
	}
	if !sp.nameSet && !sp.validNameSet {
		return nil, fmt.Errorf("must use either WithName or WithValidName")
	}

	if !isValidPackName(sp.name) {
		return nil, fmt.Errorf("invalid name: %q", sp.name)
	}

	return sp, nil
}

func (pack *StickerPack) Create() (string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	IDstring := strconv.FormatInt(pack.userID, 10)
	if err := writer.WriteField("user_id", IDstring); err != nil {
		return "", fmt.Errorf("failed to write user_id: %w", err)
	}
	if err := writer.WriteField("name", pack.name); err != nil {
		return "", fmt.Errorf("failed to write name: %w", err)
	}
	if err := writer.WriteField("title", pack.title); err != nil {
		return "", fmt.Errorf("failed to write title: %w", err)
	}

	inputStickers := make([]attachSticker, len(pack.stickers))
	for i, sticker := range pack.stickers {
		inputStickers[i] = attachSticker{
			Sticker:   fmt.Sprintf("attach://sticker%d", i),
			EmojiList: sticker.EmojiList,
			Format:    sticker.Format,
			Keywords:  sticker.Keywords,
		}
	}

	jsonStickers, err := json.Marshal(inputStickers)
	if err != nil {
		return "", fmt.Errorf("failed to convert to JSON: %w", err)
	}
	writer.WriteField("stickers", string(jsonStickers))

	for i, sticker := range pack.stickers {
		extension := ".png"
		if sticker.Format == "video" {
			extension = ".webm"
		}
		part, err := writer.CreateFormFile(fmt.Sprintf("sticker%d", i), fmt.Sprintf("sticker%d%s", i, extension))
		if err != nil {
			return "", fmt.Errorf("failed writing to request: %w", err)
		}
		if _, err := part.Write(sticker.Sticker); err != nil {
			return "", fmt.Errorf("failed attaching image: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	reqURL := requestURL("createNewStickerset")
	resp, err := http.Post(reqURL, writer.FormDataContentType(), &buf)
	if err != nil {
		return "", fmt.Errorf("createNewStickerSet failed: %w", err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("telegram API error: %s", string(body))
	}

	stickerPackURL := fmt.Sprintf("https://t.me/addstickers/%s", pack.name)
	return stickerPackURL, nil
}

func (pack *StickerPack) Delete() error {
	reqURL := requestURL("deleteStickerSet")
	resp, err := http.PostForm(reqURL, url.Values{
		"name": {pack.name},
	})
	if err != nil {
		return fmt.Errorf("deleteStickerSet failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error: %s", string(body))
	}

	return nil
}

func PackInfo(packName string) (*StickerSet, error) {
	resp, err := http.PostForm(requestURL("getStickerSet"), url.Values{
		"name": {packName},
	})
	if err != nil {
		return nil, fmt.Errorf("getStickerSet failed: %w", err)
	}
	defer resp.Body.Close()

	var result GetStickerSetResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Ok {
		return nil, fmt.Errorf("thumbnail telegram API error: %s", result.Description)
	}

	return &result.Result, nil
}

func PackThumbnailID(stickerSet *StickerSet) (string, error) {
	// thumbnail set explicitly
	if stickerSet.Thumbnail != nil && stickerSet.Thumbnail.FileID != "" {
		return stickerSet.Thumbnail.FileID, nil
	}

	// default to first sticker
	if len(stickerSet.Stickers) > 0 {
		firstSticker := stickerSet.Stickers[0]
		// sticker has a thumbnail
		// if firstSticker.Thumbnail != nil && firstSticker.Thumbnail.FileID != "" {
		// 	return firstSticker.Thumbnail.FileID, nil
		// }
		// fallback - the sticker itself

		// the thumbnail is not animated, always use the sticker
		return firstSticker.FileID, nil
	}

	return "", errors.New("no thumbnail available")
}

func (pack *StickerPack) UpdateThumbnailID() error {
	stickerSet, err := PackInfo(pack.name)
	if err != nil {
		return err
	}
	fileID, err := PackThumbnailID(stickerSet)
	if err != nil {
		return err
	}

	pack.thumbnailID = fileID
	return nil
}

func (pack *StickerPack) AddSticker(sticker InputSticker) error {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	IDstring := strconv.FormatInt(pack.userID, 10)
	if err := writer.WriteField("user_id", IDstring); err != nil {
		return fmt.Errorf("failed to write user_id: %w", err)
	}
	if err := writer.WriteField("name", pack.name); err != nil {
		return fmt.Errorf("failed to write name: %w", err)
	}

	inputSticker := attachSticker{
		Sticker:   "attach://sticker0",
		EmojiList: sticker.EmojiList,
		Format:    sticker.Format,
		Keywords:  sticker.Keywords,
	}

	jsonSticker, err := json.Marshal(inputSticker)
	if err != nil {
		return fmt.Errorf("failed to convert to JSON: %w", err)
	}
	writer.WriteField("sticker", string(jsonSticker))

	extension := ".png"
	if sticker.Format == "video" {
		extension = ".webm"
	}
	part, err := writer.CreateFormFile(
		"sticker0",
		fmt.Sprintf("sticker0%s", extension),
	)
	if err != nil {
		return fmt.Errorf("failed writing to request: %w", err)
	}
	if _, err := part.Write(sticker.Sticker); err != nil {
		return fmt.Errorf("failed attaching image: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	reqURL := requestURL("addStickerToSet")
	resp, err := http.Post(reqURL, writer.FormDataContentType(), &buf)
	if err != nil {
		return fmt.Errorf("addStickerToSet failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error: %s", string(body))
	}

	return nil
}

func (pack *StickerPack) SetTitle(title string) error {
	data := url.Values{}
	data.Set("name", pack.name)
	data.Set("title", title)

	reqURL := requestURL("setStickerSetTitle")

	resp, err := http.PostForm(reqURL, data)
	if err != nil {
		return fmt.Errorf("setStickerSetTitle failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error: %s", body)
	}

	return nil
}

func SetStickerEmojis(fileID string, emojis []string) error {
	emojiJSON, err := json.Marshal(emojis)
	if err != nil {
		return fmt.Errorf("failed to encode emoji_list: %w", err)
	}

	data := url.Values{}
	data.Set("sticker", fileID)
	data.Set("emoji_list", string(emojiJSON))

	reqURL := requestURL("setStickerEmojiList")

	resp, err := http.PostForm(reqURL, data)
	if err != nil {
		return fmt.Errorf("setStickerEmojiList failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error: %s", body)
	}

	return nil
}

func DeleteSticker(fileID string) error {
	data := url.Values{}
	data.Set("sticker", fileID)

	reqURL := requestURL("deleteStickerFromSet")

	resp, err := http.PostForm(reqURL, data)
	if err != nil {
		return fmt.Errorf("deleteStickerFromSet failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error: %s", body)
	}

	return nil
}

func SetStickerPosition(fileID string, position int) error {
	data := url.Values{}
	data.Set("sticker", fileID)
	data.Set("position", strconv.Itoa(position))

	reqURL := requestURL("setStickerPositionInSet")

	resp, err := http.PostForm(reqURL, data)
	if err != nil {
		return fmt.Errorf("setStickerPositionInSet failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error: %s", body)
	}

	return nil
}

func (pack *StickerPack) Fetch(ctx context.Context) (*StickerSet, error) {
	return FetchPack(ctx, pack.name)
}

func FetchPack(ctx context.Context, name string) (*StickerSet, error) {
	url := requestURL(fmt.Sprintf("getStickerSet?name=%s", name))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request: %w", err)
	}
	params := &retrier.RetryParams{
		Request: req,
		Client:  httpClient,
		Retries: fetchRetires,
	}
	resp, err := retrier.RequestWithCallback(ctx, params, fetchCallback)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseFetchResponse(resp.Body)
}

func FetchPackPreview(ctx context.Context, name string) (*PackPreview, error) {
	pack, err := FetchPack(ctx, name)
	if err != nil {
		return nil, err
	}

	thumbnailID, err := PackThumbnailID(pack)
	if err != nil {
		return nil, err
	}
	return &PackPreview{
		Title:       pack.Title,
		Name:        pack.Name,
		ThumbnailID: thumbnailID,
	}, nil
}

func fetchCallback(resp *http.Response) (bool, error) {
	switch resp.StatusCode {
	case http.StatusOK:
		return false, nil
	case http.StatusBadRequest:
		// the description says "STICKERSET_INVALID"
		// the error messages are very clear, very demure
		return false, ErrorPackNotFound
	default:
		return true, fmt.Errorf(
			"tg request for fetch pack status code %d", resp.StatusCode,
		)
	}
}

func parseFetchResponse(r io.Reader) (*StickerSet, error) {
	var set GetStickerSetResponse

	if err := json.NewDecoder(r).Decode(&set); err != nil {
		errMsg := fmt.Errorf("failed to decode set Fetch response: %w", err)
		fmt.Printf("%v\n", errMsg)
		return nil, errMsg
	}
	return &set.Result, nil
}

func requestURL(method string) string {
	token := config.Load().TelegramToken()
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, method)
}

func isValidPackName(name string) bool {
	// English letters and digits, underscores
	// <= 64 characters
	// no consecutive underscores
	re := regexp.MustCompile(`(?i)^[a-z][a-z0-9_]*$`)
	return len(name) <= 64 &&
		re.MatchString(name) &&
		!strings.Contains(name, "__")
}
