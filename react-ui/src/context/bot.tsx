import { createContext } from "react";

export type BotContextType = {
    id: number
    name: string
};

export const BotContext = createContext<BotContextType>({
    id: 0,
    name: ""
});