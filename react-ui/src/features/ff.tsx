import { useContext } from "react";
import { BotContext } from "../context/bot";

export function QuotesUploadScreen() {
    const bot = useContext(BotContext);

    return <form method="POST" action="/api/ff/quotes" encType="multipart/form-data">
        <h3>Загрузка цитат</h3>
        <input name="bot-id" type="hidden" value={bot.id} />
        <input name="quotes" type="file" accept="*.json" />
        <input value="Отправить" name="submit" type="submit" />
    </form>
}

export function MediaUploadScreen() {
    const bot = useContext(BotContext);

    return <form action="/api/ff/media" encType="multipart/form-data" method="post">
        <input name="bot-id" type="hidden" value={bot.id} />
        <span> ID чата </span>
        <input type="number" name="chat-id" />
        <div>
            <span> Название файла (или текст для отправки) </span>
            <br />
            <textarea name="filename" rows={5} cols={50}></textarea>
        </div>
        <div>
            <span> Тип сообщения </span>
            <select name="type">
                <option value="audio">audio</option>
                <option value="video">video</option>
                <option value="image">image</option>
                <option value="text">text</option>
            </select>
        </div>
        <input type="file" name="file-data" />
        <br />
        <input value="Отправить" name="submit" type="submit" />
    </form>
}