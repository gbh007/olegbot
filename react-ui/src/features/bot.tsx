import { useEffect, useState } from "react";
import { BotModel, useBotCreate, useBotDelete, useBotGet, useBotList, useBotUpdate } from "../apiclient/api-bot";
import { BotEditorWidget, BotListWidget, BotSelectWidget } from "../widgets/bot";
import { ErrorWidget } from "../widgets/simple-component";
import { useParams } from "react-router-dom";

export function BotListScreen() {
    const [botListResponse, fetchBotList] = useBotList()
    const [botDeleteResponse, doDeleteBot] = useBotDelete()

    useEffect(() => {
        fetchBotList()
    }, [fetchBotList])

    return <div>
        <ErrorWidget value={botListResponse} />
        <ErrorWidget value={botDeleteResponse} />
        {botListResponse.data ?
            <BotListWidget
                value={botListResponse.data!}
                deleteCallback={(v: number) => {
                    if (!confirm(`Удалить бота №${v}`)) {
                        return;
                    }

                    doDeleteBot({ id: v }).then(() => {
                        fetchBotList()
                    })
                }}
            /> : null}
    </div>
}


export function BotSelectScreen(props: {
    selectCallback: (id: number, name: string) => void
}) {
    const [botListResponse, fetchBotList] = useBotList()

    useEffect(() => {
        fetchBotList()
    }, [fetchBotList])

    return <div>
        <ErrorWidget value={botListResponse} />
        {botListResponse.data ?
            <BotSelectWidget
                value={botListResponse.data!}
                selectCallback={props.selectCallback}
            /> : null}
    </div>
}

export function BotEditorScreen() {
    const params = useParams()
    const id = params.botID && params.botID != "new" ? parseInt(params.botID) : null


    const [botInfoResponse, fetchBotInfo] = useBotGet()
    const [createBotResponse, doCreateBot] = useBotCreate()
    const [updateBotResponse, doUpdateBot] = useBotUpdate()

    const [botInfo, setBotInfo] = useState<BotModel>({
        id: 0,
        create_at: new Date().toJSON(),
        enabled: false,
        name: "",
        tag: "",
        token: "",
    })

    useEffect(() => {
        if (botInfoResponse.data) {
            setBotInfo(botInfoResponse.data)
        }
    }, [botInfoResponse.data])

    if (id) {
        useEffect(() => {
            fetchBotInfo({ id: id! })
        }, [fetchBotInfo, id])
    }

    return <div>
        <ErrorWidget value={botInfoResponse} />
        <ErrorWidget value={createBotResponse} />
        <ErrorWidget value={updateBotResponse} />
        <BotEditorWidget
            value={botInfo}
            onChange={setBotInfo}
            onSave={() => {
                if (botInfo.id) {
                    doUpdateBot(botInfo)
                } else {
                    doCreateBot(botInfo)
                    // FIXME: тут надо зароутится на список ботов
                }
            }}
        />
    </div>
}