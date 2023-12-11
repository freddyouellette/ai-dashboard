import { createSelector, createSlice } from '@reduxjs/toolkit'
import axios from 'axios';
import { goToChatPage } from './page';

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
	fetch('http://localhost:8080/api/chats')
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
export const createChat = ({ botId }) => async dispatch => {
	let newChatData = {
		name: "New Chat",
		bot_id: botId,
	}
	axios.post('http://localhost:8080/api/chats', newChatData)
	.then(response => {
		console.log(response);
		dispatch(chatsSlice.actions.addChat(response.data));
		dispatch(goToChatPage(response.data))
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