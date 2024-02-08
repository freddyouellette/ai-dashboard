import { configureStore } from '@reduxjs/toolkit'
import botReducer from './bots'
import chatReducer from './chats'
import messageReducer from './messages'
import pageReducer from './page'

const store = configureStore({
	reducer: {
		bots: botReducer,
		page: pageReducer,
		chats: chatReducer,
		messages: messageReducer,
	}, 
}, window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__());

export default store;