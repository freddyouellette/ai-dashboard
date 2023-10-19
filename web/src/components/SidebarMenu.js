import { useSelector } from "react-redux";
import { selectSidebarStatus, SIDEBAR_STATUSES } from "../store/page";
import SidebarBotList from "./SidebarBotList";
import SidebarChatList from "./SidebarChatList";

export default function SidebarMenu() {
	const sidebarStatus = useSelector(selectSidebarStatus);
	
	switch (sidebarStatus) {
		case SIDEBAR_STATUSES.BOT_LIST:
			return <SidebarBotList />;
		case SIDEBAR_STATUSES.BOT_CHAT_LIST:
			return <SidebarChatList />;
		default:
			return null;
	}
}