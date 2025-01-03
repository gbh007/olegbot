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
    sticker_chance?: number
    gif_chance?: number
    create_at: string
    update_at?: string
}

interface botRequest {
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

export function useBotDelete(): [Response<void | null>, PostAction<botRequest>] {
    const [response, fetchData] = useAPIPost<botRequest, void>('/api/bot/delete')

    return [response, fetchData]
}

export function useBotGet(): [Response<BotModel | null>, PostAction<botRequest>] {
    const [response, fetchData] = useAPIPost<botRequest, BotModel>('/api/bot/get')

    return [response, fetchData]
}

export function useBotStart(): [Response<void | null>, PostAction<botRequest>] {
    const [response, fetchData] = useAPIPost<botRequest, void>('/api/bot/start')

    return [response, fetchData]
}

export function useBotStop(): [Response<void | null>, PostAction<botRequest>] {
    const [response, fetchData] = useAPIPost<botRequest, void>('/api/bot/stop')

    return [response, fetchData]
}

export function useBotRunningList(): [Response<Array<number> | null>, PostAction<void>] {
    const [response, fetchData] = useAPIGet<Array<number>>('/api/bot/running/list')

    return [response, fetchData]
}