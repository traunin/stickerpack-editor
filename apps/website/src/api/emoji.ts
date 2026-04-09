import Graphemer from 'graphemer'

const emojiRegex = /[\p{Extended_Pictographic}\p{Emoji_Component}]/u
const splitter = new Graphemer()

export function splitEmojis(text: string): string[] {
  return splitter.splitGraphemes(text)
    .filter(s => emojiRegex.test(s))
}
