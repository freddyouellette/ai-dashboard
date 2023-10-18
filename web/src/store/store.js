import { configureStore } from '@reduxjs/toolkit'
import botReducer from './bots'
import pageReducer from './page'

export default configureStore({
	reducer: {
		bots: botReducer,
		page: pageReducer,
	}, 
}, window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__());