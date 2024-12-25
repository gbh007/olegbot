import { PostAction, Response, useAPIGet, useAPIPost } from "./client-hooks"

export interface BotModel {
    id: number
    enabled: boolean
    emoji_list?: Array<string>
    emoji_chance?: number
    tags?: Array<string>
    name: string
    tag: string
    description?: string
    token: string
    allowed_chats?: Array<number>
    create_at: string
    update_at?: string
}

interface botDeleteRequest {
    id: number
}

export function useBotList(): [Response<Array<BotModel> | null>, PostAction<void>] {
    const [response, fetchData] = useAPIGet<Array<BotModel>>('/api/bot/list')

    return [response, fetchData]
}

export function useBotCreate(): [Response<void | null>, PostAction<BotModel>] {
    const [response, fetchData] = useAPIPost<BotModel, void>('/api/bot/create')

    return [response, fetchData]
}

export function useBotUpdate(): [Response<void | null>, PostAction<BotModel>] {
    const [response, fetchData] = useAPIPost<BotModel, void>('/api/bot/update')

    return [response, fetchData]
}

export function useBotDelete(): [Response<void | null>, PostAction<botDeleteRequest>] {
    const [response, fetchData] = useAPIPost<botDeleteRequest, void>('/api/bot/delete')

    return [response, fetchData]
}

export function useBotGet(): [Response<BotModel | null>, PostAction<botDeleteRequest>] {
    const [response, fetchData] = useAPIPost<botDeleteRequest, BotModel>('/api/bot/get')

    return [response, fetchData]
}