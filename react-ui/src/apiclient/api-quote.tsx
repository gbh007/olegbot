import { PostAction, Response, useAPIGet, useAPIPost } from "./client-hooks"

export interface QuoteModel {
    id: number
    bot_id: number
    text: string
    creator_id?: number
    created_in_chat_id?: number
    create_at: string
}

interface quoteRequest {
    id: number
}

interface quoteListRequest {
    bot_id: number
}

interface quoteUpdateRequest {
    id: number
    text: string
}

interface quoteCreateRequest {
    bot_id: number
    text: string
}

export function useQuoteList(): [Response<Array<QuoteModel> | null>, PostAction<quoteListRequest>] {
    const [response, fetchData] = useAPIPost<quoteListRequest, Array<QuoteModel>>('/api/quote/list')

    return [response, fetchData]
}

export function useQuoteCreate(): [Response<void | null>, PostAction<quoteCreateRequest>] {
    const [response, fetchData] = useAPIPost<quoteCreateRequest, void>('/api/quote/create')

    return [response, fetchData]
}

export function useQuoteUpdate(): [Response<void | null>, PostAction<quoteUpdateRequest>] {
    const [response, fetchData] = useAPIPost<quoteUpdateRequest, void>('/api/quote/update')

    return [response, fetchData]
}

export function useQuoteDelete(): [Response<void | null>, PostAction<quoteRequest>] {
    const [response, fetchData] = useAPIPost<quoteRequest, void>('/api/quote/delete')

    return [response, fetchData]
}

export function useQuoteGet(): [Response<QuoteModel | null>, PostAction<quoteRequest>] {
    const [response, fetchData] = useAPIPost<quoteRequest, QuoteModel>('/api/quote/get')

    return [response, fetchData]
}