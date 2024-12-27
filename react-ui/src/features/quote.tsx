import { Link, useParams } from "react-router-dom";
import { QuoteModel, useQuoteCreate, useQuoteDelete, useQuoteGet, useQuoteList, useQuoteUpdate } from "../apiclient/api-quote";
import { useContext, useEffect, useState } from "react";
import { ErrorWidget } from "../widgets/simple-component";
import { BotContext } from "../context/bot";

export function QuoteListScreen() {
    const bot = useContext(BotContext);
    const [quoteListResponse, fetchQuoteList] = useQuoteList()
    const [quoteDeleteResponse, doDeleteQuote] = useQuoteDelete()

    useEffect(() => {
        fetchQuoteList({ bot_id: bot.id })
    }, [fetchQuoteList, bot])

    return <div>
        <ErrorWidget value={quoteListResponse} />
        <ErrorWidget value={quoteDeleteResponse} />
        {quoteListResponse.data ?
            <QuoteListWidget
                value={quoteListResponse.data!}
                deleteCallback={(v: number) => {
                    if (!confirm(`Удалить цитату №${v}`)) {
                        return;
                    }

                    doDeleteQuote({ id: v }).then(() => {
                        fetchQuoteList({ bot_id: bot.id })
                    })
                }}
            /> : null}
    </div>
}

export function QuoteEditorScreen() {
    const bot = useContext(BotContext);
    const params = useParams()
    const id = params.quoteID && params.quoteID != "new" ? parseInt(params.quoteID) : null


    const [quoteInfoResponse, fetchQuoteInfo] = useQuoteGet()
    const [createQuoteResponse, doCreateQuote] = useQuoteCreate()
    const [updateQuoteResponse, doUpdateQuote] = useQuoteUpdate()

    const [quoteInfo, setQuoteInfo] = useState<QuoteModel>({
        id: 0,
        create_at: new Date().toJSON(),
        bot_id: bot.id,
        text: "",
    })

    useEffect(() => {
        if (quoteInfoResponse.data) {
            setQuoteInfo(quoteInfoResponse.data)
        }
    }, [quoteInfoResponse.data])

    useEffect(() => {
        setQuoteInfo({ ...quoteInfo, bot_id: bot.id })
    }, [bot.id])

    if (id) {
        useEffect(() => {
            fetchQuoteInfo({ id: id! })
        }, [fetchQuoteInfo, id])
    }

    return <div>
        <ErrorWidget value={quoteInfoResponse} />
        <ErrorWidget value={createQuoteResponse} />
        <ErrorWidget value={updateQuoteResponse} />
        <QuoteEditorWidget
            value={quoteInfo}
            onChange={setQuoteInfo}
            onSave={() => {
                if (quoteInfo.id) {
                    doUpdateQuote(quoteInfo)
                } else {
                    doCreateQuote(quoteInfo)
                    // FIXME: тут надо зароутится на список цитат
                }
            }}
        />
    </div>
}

function QuoteEditorWidget(props: {
    value: QuoteModel
    onChange: (v: QuoteModel) => void
    onSave: () => void
}) {
    return <div>
        <span>ID:</span>{props.value.id}<br />
        <textarea rows={5} cols={50} value={props.value.text} onChange={(e) => {
            props.onChange({ ...props.value, text: e.target.value })
        }}></textarea><br />
        <button onClick={props.onSave}>Сохранить</button>
    </div>
}

function QuoteListWidget(props: {
    value: Array<QuoteModel>
    deleteCallback: (v: number) => void
}) {
    return <table style={{ borderSpacing: "10px" }}>
        <thead>
            <tr>
                <td>ID <Link to={"/quote/edit/new"}>новый</Link></td>
                <td>Создана</td>
                <td>Текст</td>
                <td>Пользователь</td>
                <td>Чат</td>
                <td>Действия</td>
            </tr>
        </thead>
        <tbody>
            {props.value.map(quote =>
                <tr key={quote.id}>
                    <td>{quote.id}</td>
                    <td>{quote.create_at}</td>
                    <td>{quote.text}</td>
                    <td>{quote.creator_id}</td>
                    <td>{quote.created_in_chat_id}</td>
                    <td>
                        <Link to={"/quote/edit/" + quote.id}>редактировать</Link><br />
                        <button onClick={() => { props.deleteCallback(quote.id) }}>удалить</button>
                    </td>
                </tr>
            )}
        </tbody>
    </table>
}