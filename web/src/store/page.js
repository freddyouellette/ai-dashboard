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
		botSelected: null,
		chatSelected: null,
		botToUpdate: null,
	},
	reducers: {
		goToSidebarBotChatList: (state, action) => {
			state.botSelected = action.payload;
			state.status = PAGE_STATUSES.BOT_CHAT;
		},
		goToSidebarBotList: (state, action) => {
			state.status = PAGE_STATUSES.BOT_LIST;
			state.botSelected = null;
		},
		goToBotChat: (state, action) => {
			state.status = PAGE_STATUSES.BOT_CHAT;
			state.chatSelected = action.payload;
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

export const { setBotSelected, setPageStatus, goToBotEdit } = pageSlice.actions;
export const selectPageStatus = state => state.page.status;
export const selectSidebarStatus = state => state.page.sidebarStatus;
export const selectBotToUpdate = state => state.page.botToUpdate;
export default pageSlice.reducer;