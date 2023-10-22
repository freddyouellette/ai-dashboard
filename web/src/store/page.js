import { createSlice } from '@reduxjs/toolkit'

export const SIDEBAR_STATUSES = {
	BOT_LIST: "bot-list",
	BOT_CHAT_LIST: "bot-chat-list",
}

export const PAGE_STATUSES = {
	CREATE_BOT: "create-bot",
	BOT_CHAT: "bot-chat",
};

const pageSlice = createSlice({
	name: 'page',
	initialState: {
		status: null,
		sidebarStatus: SIDEBAR_STATUSES.BOT_LIST,
		selectedBot: null,
		selectedChat: null,
		botToUpdate: null,
	},
	reducers: {
		goToSidebarBotChatList: (state, action) => {
			state.selectedBot = action.payload;
			state.sidebarStatus = SIDEBAR_STATUSES.BOT_CHAT_LIST;
		},
		goToSidebarBotList: (state, action) => {
			state.sidebarStatus = SIDEBAR_STATUSES.BOT_LIST;
			state.selectedBot = null;
		},
		goToBotChat: (state, action) => {
			state.status = PAGE_STATUSES.BOT_CHAT;
			state.selectedChat = action.payload;
		},
		setPageStatus: (state, action) => {
			state.status = action.payload;
		},
		goToBotEdit: (state, action) => {
			state.botToUpdate = action.payload;
			state.status = PAGE_STATUSES.CREATE_BOT;
		},
	}
});

export const { goToSidebarBotChatList, goToSidebarBotList, setPageStatus, goToBotEdit, goToBotChat } = pageSlice.actions;
export const selectPageStatus = state => state.page.status;
export const selectSidebarStatus = state => state.page.sidebarStatus;
export const selectBotToUpdate = state => state.page.botToUpdate;
export const selectSelectedBot = state => state.page.selectedBot;
export const selectSelectedChat = state => state.page.selectedChat;
export default pageSlice.reducer;