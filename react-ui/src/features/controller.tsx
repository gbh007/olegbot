import { useEffect } from "react";
import { useBotList, useBotRunningList, useBotStart, useBotStop } from "../apiclient/api-bot";
import { ErrorWidget } from "../widgets/simple-component";

export function BotControllerScreen() {
    const [botListResponse, fetchBotList] = useBotList()
    const [botRunningListResponse, fetchBotRunningList] = useBotRunningList()

    const [botStartResponse, doStartBot] = useBotStart()
    const [botStopResponse, doStopBot] = useBotStop()

    useEffect(() => {
        fetchBotList()
    }, [fetchBotList])

    useEffect(() => {
        fetchBotRunningList()
    }, [fetchBotRunningList])

    return <div>
        <ErrorWidget value={botListResponse} />
        <ErrorWidget value={botRunningListResponse} />
        <ErrorWidget value={botStartResponse} />
        <ErrorWidget value={botStopResponse} />
        <table style={{ borderSpacing: "10px" }}>
            <thead>
                <tr>
                    <td>ID</td>
                    <td>Название</td>
                    <td>Статус</td>
                    <td>Действия</td>
                </tr>
            </thead>
            <tbody>
                {botListResponse.data?.map(bot =>
                    <tr key={bot.id}>
                        <td>{bot.id}</td>
                        <td>{bot.name}</td>
                        <td>{botRunningListResponse.data?.includes(bot.id) ? 'Запущен' : bot.enabled ? 'Готов к запуску' : 'Не включен'}</td>
                        <td>
                            <button onClick={() => {
                                doStartBot({ id: bot.id }).then((() => { fetchBotRunningList() }))
                            }} disabled={botStartResponse.isLoading}>Запустить</button><br />
                            <button onClick={() => {
                                doStopBot({ id: bot.id }).then((() => { fetchBotRunningList() }))
                            }} disabled={botStopResponse.isLoading}>Остановить</button>
                        </td>
                    </tr>
                )}
            </tbody>
        </table>
    </div>
}