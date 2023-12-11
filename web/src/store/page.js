import { createSlice } from '@reduxjs/toolkit'

export const PAGE_STATUSES = {
	CREATE_BOT: "create-bot",
	BOT_LIST: "bot-list",
	CREATE_CHAT: "create-chat",
	CHAT_LIST: "chat-list",
	BOT_CHAT: "bot-chat",
};

const pageSlice = createSlice({
	name: 'page',
	initialState: {
		status: null,
		selectedChat: null,
		selectedBot: null,
	},
	reducers: {
		setStatus: (state, action) => {
			state.status = action.payload;
		},
		setSelectedChat: (state, action) => {
			state.selectedChat = action.payload;
		},
		setSelectedBot: (state, action) => {
			state.selectedBot = action.payload;
		},
	}
});

export const goToCreateChatPage = () => dispatch => {
	dispatch(pageSlice.actions.setStatus(PAGE_STATUSES.CREATE_CHAT))
}

export const goToBotListPage = () => dispatch => {
	dispatch(pageSlice.actions.setStatus(PAGE_STATUSES.BOT_LIST))
}

export const goToChatListPage = () => dispatch => {
	dispatch(pageSlice.actions.setStatus(PAGE_STATUSES.CHAT_LIST))
}

// thunk
export const goToChatPage = (chat) => (dispatch, getState) => {
	dispatch(pageSlice.actions.setSelectedChat(chat))
	dispatch(pageSlice.actions.setStatus(PAGE_STATUSES.BOT_CHAT))
}

// thunk
export const goToBotEditPage = (bot) => (dispatch, getState) => {
	dispatch(pageSlice.actions.setStatus(PAGE_STATUSES.CREATE_BOT))
	dispatch(pageSlice.actions.setSelectedBot(bot))
}

export const selectPageStatus = state => state.page.status;
export const selectSelectedChat = state => state.page.selectedChat;
export const selectSelectedBot = state => state.page.selectedBot;

export default pageSlice.reducer;