import { PropsWithChildren, useContext, useState } from 'react'
import {
  RouterProvider,
  Link,
  createHashRouter,
} from "react-router-dom";
import { BotEditorScreen, BotListScreen, BotSelectScreen } from './features/bot';
import { BotContext, BotContextType } from './context/bot';
import { QuoteEditorScreen, QuoteListScreen } from './features/quote';
import { ModeratorsScreen } from './features/moderator';
import { MediaUploadScreen, QuotesUploadScreen } from './features/ff';


const router = createHashRouter([
  {
    path: "/",
    element: <SimpleWrapper></SimpleWrapper>,
  },
  {
    path: "/bot/list",
    element: <SimpleWrapper> <BotListScreen /></SimpleWrapper>,
  },
  {
    path: "/bot/edit/:botID",
    element: <SimpleWrapper> <BotEditorScreen /></SimpleWrapper>,
  },
  {
    path: "/quote/list",
    element: <SimpleWrapper> <QuoteListScreen /></SimpleWrapper>,
  },
  {
    path: "/quote/edit/:quoteID",
    element: <SimpleWrapper> <QuoteEditorScreen /></SimpleWrapper>,
  },
  {
    path: "/moderator/list",
    element: <SimpleWrapper> <ModeratorsScreen /></SimpleWrapper>,
  },
  {
    path: "/ff/quotes",
    element: <SimpleWrapper> <QuotesUploadScreen /></SimpleWrapper>,
  },
  {
    path: "/ff/media",
    element: <SimpleWrapper> <MediaUploadScreen /></SimpleWrapper>,
  },
]);

function App() {
  const [bot, setBot] = useState<BotContextType>({
    id: 0,
    name: ""
  });

  return (
    <>
      {bot.id > 0 ?
        <BotContext.Provider value={bot}>
          <RouterProvider router={router} />
        </BotContext.Provider>
        :
        <BotSelectScreen selectCallback={(id: number, name: string) => {
          setBot({
            id: id,
            name: name
          })
        }} />
      }
    </ >
  )
}

function SimpleWrapper(props: PropsWithChildren) {
  const bot = useContext(BotContext);

  return (
    <>
      <div>
        <Link to="/">Главная </Link>
        <Link to="/bot/list">Список ботов </Link>
        <Link to="/quote/list">Список цитат </Link>
        <Link to="/moderator/list">Список модераторов </Link>
        <Link to="/ff/quotes">Загрузка цитат </Link>
        <Link to="/ff/media">Отправка медиа </Link>
        <span>Выбран бот: {bot.name}</span>
      </div>

      {props.children}
    </ >
  )
}

export default App
