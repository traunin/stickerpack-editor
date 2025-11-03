export interface PackResponse {
  id: number
  title: string
  name: string
  thumbnail_id: string
}

export interface PackParameters {
  name: string
  title: string
  hasWatermark: boolean
  isPublic: boolean
}
