import { Link } from "react-router-dom";
import { BotModel } from "../apiclient/api-bot";
import { StringArrayPicker } from "./simple-component";

export function BotListWidget(props: {
    value: Array<BotModel>
    deleteCallback: (v: number) => void
}) {
    return <table>
        <thead>
            <tr>
                <td>ID <Link to={"/bot/edit/new"}>новый</Link></td>
                <td>Название</td>
                <td>Создан</td>
                <td>Действия</td>
            </tr>
        </thead>
        <tbody>
            {props.value.map(bot =>
                <tr key={bot.id}>
                    <td>{bot.id}</td>
                    <td>{bot.name}</td>
                    <td>{bot.create_at}</td>
                    <td>
                        <Link to={"/bot/edit/" + bot.id}>редактировать</Link>
                        <button onClick={() => { props.deleteCallback(bot.id) }}>удалить</button>
                    </td>
                </tr>
            )}
        </tbody>
    </table>
}

export function BotEditorWidget(props: {
    value: BotModel
    onChange: (v: BotModel) => void
    onSave: () => void
}) {
    return <div>
        <span>ID:</span>{props.value.id}<br />
        <label>Включен: <input
            type="checkbox"
            checked={props.value.enabled}
            onChange={(e) => {
                props.onChange({ ...props.value, enabled: e.target.checked })
            }}
        /></label><br />
        <span>Название: </span><input
            type="text"
            value={props.value.name}
            onChange={(e) => {
                props.onChange({ ...props.value, name: e.target.value })
            }}
        /><br />
        <span>Тег в ТГ: </span><input
            type="text"
            value={props.value.tag}
            onChange={(e) => {
                props.onChange({ ...props.value, tag: e.target.value })
            }}
        /><br />
        <span>Описание: </span><input
            type="text"
            value={props.value.description ?? ""}
            onChange={(e) => {
                props.onChange({ ...props.value, description: e.target.value })
            }}
        /><br />
        <span>Токен: </span><input
            type="text"
            value={props.value.token}
            onChange={(e) => {
                props.onChange({ ...props.value, token: e.target.value })
            }}
        /><br />
        <StringArrayPicker
            name="Реакции"
            value={props.value.emoji_list}
            onChange={e => {
                props.onChange({ ...props.value, emoji_list: e })
            }}
        />
        <span>Шанс реакции: </span><input
            type="number"
            step={0.001}
            max={1}
            min={0}
            value={props.value.emoji_chance ?? 0}
            onChange={(e) => {
                props.onChange({ ...props.value, emoji_chance: e.target.valueAsNumber })
            }}
        /><br />
        <StringArrayPicker
            name="Теги"
            value={props.value.tags}
            onChange={e => {
                props.onChange({ ...props.value, tags: e })
            }}
        />
        <StringArrayPicker
            name="Разрешенные чаты"
            value={props.value.allowed_chats?.map(e => e.toString())}
            onChange={e => {
                props.onChange({ ...props.value, allowed_chats: e.map(e => e ? parseInt(e) : 0) })
            }}
        />
        <button onClick={props.onSave}>Сохранить</button>
    </div>
}