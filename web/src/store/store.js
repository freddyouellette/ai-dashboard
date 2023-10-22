import { configureStore } from '@reduxjs/toolkit'
import botReducer, { fetchBots } from './bots'
import chatReducer, { fetchChats } from './chats'
import messageReducer, { fetchMessages } from './messages'
import pageReducer from './page'

const store = configureStore({
	reducer: {
		bots: botReducer,
		page: pageReducer,
		chats: chatReducer,
		messages: messageReducer,
	}, 
}, window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__());

// initialize the store with the starting data from the server
store.dispatch(fetchBots())
store.dispatch(fetchChats())
store.dispatch(fetchMessages())

export default store;