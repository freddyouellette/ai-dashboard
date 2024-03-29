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

// thunk
export const goToChatEditPage = (chat) => dispatch => {
	dispatch(pageSlice.actions.setStatus(PAGE_STATUSES.CREATE_CHAT))
	dispatch(pageSlice.actions.setSelectedBot(null))
	dispatch(pageSlice.actions.setSelectedChat(chat || null))
	document.title = chat ? ("Edit: " + chat.name) : "Create Chat"
}

// thunk
export const goToBotListPage = () => dispatch => {
	dispatch(pageSlice.actions.setStatus(PAGE_STATUSES.BOT_LIST))
	dispatch(pageSlice.actions.setSelectedBot(null))
	dispatch(pageSlice.actions.setSelectedChat(null))
	document.title = "AI Dashboard"
}

// thunk
export const goToChatListPage = () => dispatch => {
	dispatch(pageSlice.actions.setStatus(PAGE_STATUSES.CHAT_LIST))
	dispatch(pageSlice.actions.setSelectedBot(null))
	dispatch(pageSlice.actions.setSelectedChat(null))
	document.title = "AI Dashboard"
}

// thunk
export const goToChatPage = (chat) => (dispatch) => {
	dispatch(pageSlice.actions.setSelectedChat(chat))
	dispatch(pageSlice.actions.setStatus(PAGE_STATUSES.BOT_CHAT))
	dispatch(pageSlice.actions.setSelectedBot(null))
	document.title = chat.name
}

// thunk
export const goToBotEditPage = (bot) => (dispatch) => {
	dispatch(pageSlice.actions.setStatus(PAGE_STATUSES.CREATE_BOT))
	dispatch(pageSlice.actions.setSelectedChat(null))
	dispatch(pageSlice.actions.setSelectedBot(bot))
	document.title = bot ? ("Edit: " + bot.name) : "Create Bot"
}

export const selectPageStatus = state => state.page.status;
export const selectSelectedChat = state => state.page.selectedChat;
export const selectSelectedBot = state => state.page.selectedBot;

export default pageSlice.reducer;