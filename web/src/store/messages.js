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
	}
});

// thunk
export const sendMessage = (chatId, botId, message) => async dispatch => {
	if (message) {
		let newMessageData = {
			chat_id: chatId,
			text: message,
		}
		await axios.post(process.env.REACT_APP_API_HOST+'/api/messages', newMessageData)
		.then(response => {
			dispatch(messagesSlice.actions.addMessage(response.data))
			dispatch(messagesSlice.actions.setWaitingForResponse(true))
			dispatch(getMessageCorrection(chatId, response.data.ID))
			// refresh bots so the order corrects itself
			dispatch(getBot(botId))
		}, error => console.error(error))
	} else {
		dispatch(messagesSlice.actions.setWaitingForResponse(true))
	}
	axios.get(process.env.REACT_APP_API_HOST+"/api/chats/"+chatId+"/response")
	.then(response => {
		dispatch(messagesSlice.actions.setWaitingForResponse(false))
		dispatch(messagesSlice.actions.addMessage(response.data))
	}, error => console.error(error))
}

// thunk
export const getMessageCorrection = (chatId, messageId) => async (dispatch, getState) => {
	let chat = getState().chats.chats[chatId]
	let bot = getState().bots.bots[chat.bot_id]
	if (bot.correction_prompt) {
		dispatch(messagesSlice.actions.setWaitingForCorrectionId(messageId))
		axios.get(process.env.REACT_APP_API_HOST+"/api/messages/"+messageId+"/correction")
		.then(response => {
			dispatch(messagesSlice.actions.addMessage(response.data))
			dispatch(messagesSlice.actions.setWaitingForCorrectionId(null))
		}, error => console.error(error))
	}
}

// thunk
export const getChatMessages = (chat, page) => async dispatch => {
	dispatch(messagesSlice.actions.setMessagesLoading(true));
	dispatch(messagesSlice.actions.setMessagesError(null));
	let params = {
		chat_id: chat.ID,
		page,
	};
	axios.get(process.env.REACT_APP_API_HOST+`/api/messages`, { params })
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
	(messages, messagesLoading, messagesError, messagesPage, messagesTotal) => ({ messages, messagesLoading, messagesError, messagesPage, messagesTotal })
);
export const selectWaitingForResponse = state => state.messages.waitingForResponse
export const selectWaitingForCorrectionId = state => state.messages.waitingForCorrectionId

export default messagesSlice.reducer;