import axios from "axios";
import { useCallback, useState } from "react";

export type GetAction = () => Promise<void>
export type PostAction<T> = (data: T) => Promise<void>

export interface Response<T> {
    data: T,
    isLoading: boolean,
    isError: boolean,
    isUnauthorize: boolean,
    errorText: string,
}

export function useAPIGet<ResponseType>(url: string): [Response<ResponseType | null>, GetAction] {
    const [data, setData] = useState<ResponseType | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [isError, setIsError] = useState(false);
    const [isUnauthorize, setIsUnauthorize] = useState(false);
    const [errorText, setErrorText] = useState('');


    const fetchData = useCallback(async () => {
        if (!url) {
            return
        }

        setIsError(false);
        setIsUnauthorize(false)
        setIsLoading(true);
        setErrorText('');

        try {
            const result = await axios.get(url);

            setData(result.data);
        } catch (error) {
            let text = 'unknown error'

            if (error instanceof axios.AxiosError) {
                text = error.response?.data ?? 'unknown error'

                if (error.response?.status == 401) {
                    setIsUnauthorize(true)
                }

                if (error.response?.status == 403) {
                    setIsUnauthorize(true)
                }

            }

            setErrorText(text);
            setIsError(true);
            setIsLoading(false);

            throw text
        }

        setIsLoading(false);
    }, [url])

    return [{ data, isLoading, isError, errorText, isUnauthorize }, fetchData];
}

export function useAPIPost<RequestType, ResponseType>(url: string): [Response<ResponseType | null>, PostAction<RequestType>] {
    const [data, setData] = useState<ResponseType | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [isError, setIsError] = useState(false);
    const [isUnauthorize, setIsUnauthorize] = useState(false);
    const [errorText, setErrorText] = useState('');


    const fetchData = useCallback(async (requestData: RequestType) => {
        if (!url) {
            return
        }

        setIsError(false);
        setIsUnauthorize(false)
        setIsLoading(true);
        setErrorText('');

        try {
            const result = await axios.post(url, requestData);

            setData(result.data);
        } catch (error) {
            let text = 'unknown error'

            if (error instanceof axios.AxiosError) {
                text = error.response?.data ?? 'unknown error'

                if (error.response?.status == 401) {
                    setIsUnauthorize(true)
                }

                if (error.response?.status == 403) {
                    setIsUnauthorize(true)
                }

            }

            setErrorText(text);
            setIsError(true);
            setIsLoading(false);

            throw text
        }

        setIsLoading(false);
    }, [url])

    return [{ data, isLoading, isError, errorText, isUnauthorize }, fetchData];
}
