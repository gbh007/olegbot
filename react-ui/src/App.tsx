import { PropsWithChildren } from 'react'
import {
  RouterProvider,
  Link,
  createHashRouter,
} from "react-router-dom";
import { BotEditorScreen, BotListScreen } from './features/bot';


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
]);

function App() {
  return (
    <>
      <RouterProvider router={router} />
    </ >
  )
}

function SimpleWrapper(props: PropsWithChildren) {

  return (
    <>
      <div>
        <Link to="/">Главная</Link>
        <Link to="/bot/list">Список ботов</Link>
      </div>

      {props.children}
    </ >
  )
}

export default App
