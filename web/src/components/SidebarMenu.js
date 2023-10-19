import { useSelector } from "react-redux";
import { selectSidebarStatus, SIDEBAR_STATUSES } from "../store/page";
import SidebarBotList from "./SidebarBotList";

export default function SidebarMenu() {
	const sidebarStatus = useSelector(selectSidebarStatus);
	
	switch (sidebarStatus) {
		case SIDEBAR_STATUSES.BOT_LIST:
			return <SidebarBotList />;
		// case SIDEBAR_STATUSES.BOT_CHAT_LIST:
		// 	return <SidebarBotChatList />;
		default:
			return null;
	}
}