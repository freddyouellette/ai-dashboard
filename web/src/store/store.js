import { configureStore } from '@reduxjs/toolkit'
import { botsSlice } from './bots'

export default configureStore({
  reducer: {
	bots: botsSlice.reducer
  }
})