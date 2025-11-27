import Graphemer from 'graphemer'

const splitter = new Graphemer()

export function splitEmojis(text: string): string[] {
  return splitter.splitGraphemes(text)
    .filter(s => /[\p{Extended_Pictographic}\p{Emoji_Component}]/u.test(s))
}
