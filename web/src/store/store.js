import { configureStore } from '@reduxjs/toolkit'
import botReducer, { fetchBots } from './bots'
import chatReducer from './chats'
import pageReducer from './page'

const store = configureStore({
	reducer: {
		bots: botReducer,
		page: pageReducer,
		chats: chatReducer,
	}, 
}, window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__());

// initialize the store with the starting data from the server
store.dispatch(fetchBots())

export default store;