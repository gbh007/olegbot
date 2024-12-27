import { PostAction, Response, useAPIPost } from "./client-hooks"

export interface ModeratorModel {
    bot_id: number
    user_id: number
    description?: string
    create_at: string
}

interface moderatorDeleteRequest {
    bot_id: number
    user_id: number
}

interface moderatorListRequest {
    bot_id: number
}

interface moderatorCreateRequest {
    bot_id: number
    user_id: number
    description?: string
}

export function useModeratorList(): [Response<Array<ModeratorModel> | null>, PostAction<moderatorListRequest>] {
    const [response, fetchData] = useAPIPost<moderatorListRequest, Array<ModeratorModel>>('/api/moderator/list')

    return [response, fetchData]
}

export function useModeratorCreate(): [Response<void | null>, PostAction<moderatorCreateRequest>] {
    const [response, fetchData] = useAPIPost<moderatorCreateRequest, void>('/api/moderator/create')

    return [response, fetchData]
}

export function useModeratorDelete(): [Response<void | null>, PostAction<moderatorDeleteRequest>] {
    const [response, fetchData] = useAPIPost<moderatorDeleteRequest, void>('/api/moderator/delete')

    return [response, fetchData]
}