import { createSelector, createSlice } from '@reduxjs/toolkit'
import axios from 'axios';

const messagesSlice = createSlice({
	name: 'messages',
	initialState: {
		messages: {},
		messagesLoading: false,
		messagesError: null,
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
		setWaitingForCorrectionId: (state, action) => {
			state.waitingForCorrectionId = action.payload
		},
		setWaitingForResponse: (state, action) => {
			state.waitingForResponse = action.payload
		},
	}
});

// thunk
export const sendMessage = (chatId, message) => async dispatch => {
	if (message) {
		let newMessageData = {
			chat_id: chatId,
			text: message,
		}
		await axios.post(process.env.REACT_APP_API_HOST+'/api/messages', newMessageData)
		.then(response => {
			console.log(response)
			dispatch(messagesSlice.actions.addMessage(response.data))
			dispatch(messagesSlice.actions.setWaitingForResponse(true))
			dispatch(getMessageCorrection(chatId, response.data.ID))
		}, error => console.error(error))
	} else {
		dispatch(messagesSlice.actions.setWaitingForResponse(true))
	}
	axios.get(process.env.REACT_APP_API_HOST+"/api/chats/"+chatId+"/response")
	.then(response => {
		console.log(response)
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
			console.log(response)
			dispatch(messagesSlice.actions.addMessage(response.data))
			dispatch(messagesSlice.actions.setWaitingForCorrectionId(null))
		}, error => console.error(error))
	}
}

// thunk
export const getChatMessages = chat => async dispatch => {
	dispatch(messagesSlice.actions.setMessagesLoading(true));
	dispatch(messagesSlice.actions.setMessagesError(null));
	axios.get(process.env.REACT_APP_API_HOST+`/api/chats/${chat.ID}/messages`)
	.then(response => {
		console.log(response)
		dispatch(messagesSlice.actions.setMessagesLoading(false));
		dispatch(messagesSlice.actions.setMessagesError(null));
		let messages = {};
		for (let message of response.data) {
			messages[message.ID] = message
		}
		dispatch(messagesSlice.actions.setMessages(messages));
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
	(messages, messagesLoading, messagesError) => ({ messages, messagesLoading, messagesError })
);
export const selectWaitingForResponse = state => state.messages.waitingForResponse
export const selectWaitingForCorrectionId = state => state.messages.waitingForCorrectionId

export default messagesSlice.reducer;