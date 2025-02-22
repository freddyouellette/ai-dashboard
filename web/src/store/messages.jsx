import { createSelector, createSlice } from '@reduxjs/toolkit'
import axios from 'axios';
import { getBot } from './bots';

const messagesSlice = createSlice({
	name: 'messages',
	initialState: {
		messages: {},
		messagesLoading: false,
		messagesError: null,
		messagesPage: 1,
		messagesTotal: 0,
		waitingForCorrectionId: null,
		waitingForResponse: false,
		responseFailed: false,
	},
	reducers: {
		addMessage: (state, action) => {
			state.messages[action.payload.ID] = action.payload
		},
		setMessages: (state, action) => {
			state.messages = action.payload
		},
		setMessagesLoading: (state, action) => {
			state.messagesLoading = action.payload
		},
		setMessagesError: (state, action) => {
			state.messagesError = action.payload
		},
		setMessagesPage: (state, action) => {
			state.messagesPage = action.payload
		},
		setMessagesTotal: (state, action) => {
			state.messagesTotal = action.payload
		},
		setWaitingForCorrectionId: (state, action) => {
			state.waitingForCorrectionId = action.payload
		},
		setWaitingForResponse: (state, action) => {
			state.waitingForResponse = action.payload
		},
		setResponseFailed: (state, action) => {
			state.responseFailed = action.payload
		},
	}
});

export const updateMessage = (message) => async dispatch => {
	await axios.put(import.meta.env.VITE_API_HOST+'/api/messages/'+message.ID, message)
	.then(response => {
		dispatch(messagesSlice.actions.addMessage(response.data))
	}, error => console.error(error))
}

// thunk
export const sendMessage = (chatId, botId, message, role, getResponse) => async dispatch => {
	dispatch(messagesSlice.actions.setResponseFailed(false))
	if (message) {
		let newMessageData = {
			chat_id: chatId,
			text: message,
			role: role,
		}
		await axios.post(import.meta.env.VITE_API_HOST+'/api/messages', newMessageData)
		.then(response => {
			dispatch(messagesSlice.actions.addMessage(response.data))
			dispatch(getMessageCorrection(chatId, response.data.ID))
			// refresh bots so the order corrects itself
			dispatch(getBot(botId))
		}, error => {
			console.error(error)
			dispatch(messagesSlice.actions.setResponseFailed(true))
		})
	}
	if (getResponse) {
		dispatch(messagesSlice.actions.setWaitingForResponse(true))
		await axios.get(import.meta.env.VITE_API_HOST+"/api/chats/"+chatId+"/response")
		.then(response => {
			dispatch(messagesSlice.actions.setWaitingForResponse(false))
			dispatch(messagesSlice.actions.addMessage(response.data))
		}, error => {
			console.error(error)
			dispatch(messagesSlice.actions.setWaitingForResponse(false))
			dispatch(messagesSlice.actions.setResponseFailed(true))
		})
	}
}

// thunk
export const getMessageCorrection = (chatId, messageId) => async (dispatch, getState) => {
	let chat = getState().chats.chats[chatId]
	let bot = getState().bots.bots[chat.bot_id]
	if (bot.correction_prompt) {
		dispatch(messagesSlice.actions.setWaitingForCorrectionId(messageId))
		axios.get(import.meta.env.VITE_API_HOST+"/api/messages/"+messageId+"/correction")
		.then(response => {
			dispatch(messagesSlice.actions.addMessage(response.data))
			dispatch(messagesSlice.actions.setWaitingForCorrectionId(null))
		}, error => console.error(error))
	}
}

// thunk
export const updateMessageActive = (message, active) => async (dispatch, getState) => {
	const updatedMessage = { ...message, active };
	dispatch(updateMessage(updatedMessage));
}

// thunk
export const updateMessageBreakAfter = (message, break_after) => async (dispatch, getState) => {
	const updatedMessage = { ...message, break_after };
	dispatch(updateMessage(updatedMessage));
}

// thunk
export const getChatMessages = (chat, page) => async dispatch => {
	dispatch(messagesSlice.actions.setMessagesLoading(true));
	dispatch(messagesSlice.actions.setMessagesError(null));
	dispatch(messagesSlice.actions.setResponseFailed(false))
	let params = {
		chat_id: chat.ID,
		page,
	};
	axios.get(import.meta.env.VITE_API_HOST+`/api/messages`, { params })
	.then(response => {
		dispatch(messagesSlice.actions.setMessagesLoading(false));
		dispatch(messagesSlice.actions.setMessagesError(null));
		let messages = {};
		for (let message of response.data.messages) {
			messages[message.ID] = message
		}
		dispatch(messagesSlice.actions.setMessages(messages));
		dispatch(messagesSlice.actions.setMessagesPage(response.data.page));
		dispatch(messagesSlice.actions.setMessagesTotal(response.data.total));
	}, error => {
		dispatch(messagesSlice.actions.setMessagesLoading(false));
		dispatch(messagesSlice.actions.setMessagesError(error));
		console.error(error)
	});
}

export const selectMessages = createSelector(
	state => state.messages.messages,
	state => state.messages.messagesLoading,
	state => state.messages.messagesError,
	state => state.messages.messagesPage,
	state => state.messages.messagesTotal,
	state => state.messages.responseFailed,
	(messages, messagesLoading, messagesError, messagesPage, messagesTotal, responseFailed) => ({ messages, messagesLoading, messagesError, messagesPage, messagesTotal, responseFailed })
);
export const selectWaitingForResponse = state => state.messages.waitingForResponse
export const selectWaitingForCorrectionId = state => state.messages.waitingForCorrectionId

export default messagesSlice.reducer;