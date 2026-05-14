import GraphemerImport from 'graphemer'

interface GraphemerModule { default: typeof GraphemerImport }
const Graphemer =
  (GraphemerImport as unknown as GraphemerModule).default ?? GraphemerImport
const emojiRegex = /[\p{Extended_Pictographic}\p{Emoji_Component}]/u
const splitter = new Graphemer()

export function splitEmojis(text: string): string[] {
  return splitter.splitGraphemes(text)
    .filter((s: string) => emojiRegex.test(s))
}
