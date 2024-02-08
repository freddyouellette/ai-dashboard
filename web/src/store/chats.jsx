import { createSelector, createSlice } from '@reduxjs/toolkit'
import axios from 'axios';

const chatsSlice = createSlice({
	name: 'chats',
	initialState: {
		chats: {},
		chatsLoading: false,
		chatsError: null,
		selectedChat: null,
	},
	reducers: {
		setChats: (state, action) => {
			state.chats = action.payload;
		},
		addChat: (state, action) => {
			state.chats[action.payload.ID] = action.payload;
		},
		setChatsLoading: (state, action) => {
			state.chatsLoading = action.payload;
		},
		setChatsError: (state, action) => {
			state.chatsError = action.payload;
		},
	}
});

// thunk
export const getChats = () => async dispatch => {
	dispatch(chatsSlice.actions.setChatsLoading(true))
	dispatch(chatsSlice.actions.setChatsError(null))
	fetch(import.meta.env.VITE_API_HOST+'/api/chats')
	.then(response => response.json())
	.then(
		chats => {
			let chatsById = {};
			chats.forEach(chat => {
				chatsById[chat.ID] = chat;
			});
			dispatch(chatsSlice.actions.setChats(chatsById))
			dispatch(chatsSlice.actions.setChatsLoading(false))
			dispatch(chatsSlice.actions.setChatsError(null))
		},
		error => {
			dispatch(chatsSlice.actions.setChatsLoading(false))
			dispatch(chatsSlice.actions.setChatsError(error))
			console.error(error)
		},
	);
}

// thunk
export const persistChat = (chat) => async dispatch => {
	chat.memory_duration = parseInt(chat.memory_duration);
	return axios[chat.ID ? "put" : "post"](import.meta.env.VITE_API_HOST+'/api/chats', chat)
	.then(response => {
		dispatch(chatsSlice.actions.addChat(response.data));
		return response.data;
	}).catch(error => {
		console.error(error);
	})
}

const selectChatsSimple = state => state.chats.chats;
const selectChatsLoading = state => state.chats.chatsLoading;
const selectChatsError = state => state.chats.chatsError;

export const selectChats = createSelector(
	selectChatsSimple, selectChatsLoading, selectChatsError, 
	(chats, chatsLoading, chatsError) => ({ chats, chatsLoading, chatsError })
);

export default chatsSlice.reducer;