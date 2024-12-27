import { useContext, useEffect } from "react";
import { BotContext } from "../context/bot";
import { useModeratorCreate, useModeratorDelete, useModeratorList } from "../apiclient/api-moderator";
import { ErrorWidget } from "../widgets/simple-component";

export function ModeratorsScreen() {
    const bot = useContext(BotContext);

    const [moderatorListResponse, fetchModeratorList] = useModeratorList()
    const [createModeratorResponse, doCreateModerator] = useModeratorCreate()
    const [deleteModeratorResponse, doDeleteModerator] = useModeratorDelete()

    useEffect(() => {
        fetchModeratorList({ bot_id: bot.id })
    }, [fetchModeratorList, bot.id])


    return <div>
        <ErrorWidget value={moderatorListResponse} />
        <ErrorWidget value={createModeratorResponse} />
        <ErrorWidget value={deleteModeratorResponse} />
        {moderatorListResponse.data ?
            <table style={{ borderSpacing: "10px" }}>
                <thead>
                    <tr>
                        <td>
                            ID <button onClick={() => {
                                const userID = parseInt(prompt("Введите ID пользователя") || "0");
                                if (!userID) {
                                    return;
                                }

                                const text = prompt("Введите описание");

                                doCreateModerator({
                                    bot_id: bot.id,
                                    description: text || "",
                                    user_id: userID,
                                }).then(() => fetchModeratorList({ bot_id: bot.id }))
                            }}>новый</button>
                        </td>
                        <td>Добавлен</td>
                        <td>Описание</td>
                        <td>Действия</td>
                    </tr>
                </thead>
                <tbody>
                    {moderatorListResponse.data.map(moderator =>
                        <tr key={moderator.user_id}>
                            <td>{moderator.user_id}</td>
                            <td>{moderator.create_at}</td>
                            <td>{moderator.description}</td>
                            <td>
                                <button onClick={() => {
                                    if (!confirm(`Удалить модератора #${moderator.user_id}`)) {
                                        return;
                                    }

                                    doDeleteModerator({
                                        bot_id: moderator.bot_id,
                                        user_id: moderator.user_id,
                                    }).then(() => fetchModeratorList({ bot_id: bot.id }))
                                }}>удалить</button>
                            </td>
                        </tr>
                    )}
                </tbody>
            </table> : null}
    </div>
}