import { createSlice } from '@reduxjs/toolkit'

const pageSlice = createSlice({
	name: 'page',
	initialState: {
		status: null,
		botToUpdate: null,
	},
	reducers: {
		changePageStatus: (state, action) => {
			state.status = action.payload;
		},
		setBotToUpdate: (state, action) => {
			state.botToUpdate = action.payload;
			state.status = PAGE_STATUSES.CREATE_BOT;
		},
	}
});

export const PAGE_STATUSES = {
	CREATE_BOT: "create-bot",
};
export const { changePageStatus, setBotToUpdate } = pageSlice.actions;
export const selectPageStatus = state => state.page.status;
export const selectBotToUpdate = state => state.page.botToUpdate;
export default pageSlice.reducer;