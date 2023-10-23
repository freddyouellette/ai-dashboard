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
		sidebarBotSelected: null,
		selectedChat: null,
		chatBot: null,
		botToUpdate: null,
	},
	reducers: {
		goToSidebarBotChatList: (state, action) => {
			state.sidebarBotSelected = action.payload;
			state.sidebarStatus = SIDEBAR_STATUSES.BOT_CHAT_LIST;
		},
		goToSidebarBotList: (state, action) => {
			state.sidebarStatus = SIDEBAR_STATUSES.BOT_LIST;
			state.sidebarBotSelected = null;
		},
		setSelectedChat: (state, action) => {
			state.status = PAGE_STATUSES.BOT_CHAT;
			state.selectedChat = action.payload;
		},
		setChatBot: (state, action) => {
			state.chatBot = action.payload;
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

// thunk
export const goToBotChat = (chat) => (dispatch, getState) => {
	dispatch(pageSlice.actions.setSelectedChat(chat))
	dispatch(pageSlice.actions.setPageStatus(PAGE_STATUSES.BOT_CHAT))
	
	let chatBot = getState().bots.bots.find(bot => bot.ID === chat.bot_id)
	dispatch(pageSlice.actions.setChatBot(chatBot))
}

export const { goToSidebarBotChatList, goToSidebarBotList, setPageStatus, goToBotEdit } = pageSlice.actions;
export const selectPageStatus = state => state.page.status;
export const selectSidebarStatus = state => state.page.sidebarStatus;
export const selectBotToUpdate = state => state.page.botToUpdate;
export const selectSidebarSelectedBot = state => state.page.sidebarBotSelected;
export const selectSelectedChat = state => state.page.selectedChat;
export const selectChatBot = state => state.page.chatBot;
export default pageSlice.reducer;