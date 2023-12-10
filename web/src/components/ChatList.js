import { useDispatch, useSelector } from "react-redux";
import { Container, ListGroup, ListGroupItem } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faChevronRight, faPlus } from "@fortawesome/free-solid-svg-icons";
import { goToBotChat } from "../store/page";
import { getChats, selectChats } from "../store/chats";
import { TODO } from "../util/todo";
import { useEffect } from "react";
import { getBots, selectBots } from "../store/bots";

export default function ChatList() {
	const dispatch = useDispatch();
	useEffect(() => {
		dispatch(getChats())
		dispatch(getBots())
	}, [dispatch]);
	let { chats, chatsLoading, chatsError } = useSelector(selectChats);
	let { bots, botsLoading, botsError } = useSelector(selectBots);
	
	if (chatsLoading || botsLoading) return <div>Loading...</div>;
	if (chatsError) return <div>Error loading chats...</div>;
	if (botsError) return <div>Error loading bots...</div>;
	
	return (
		<ListGroup className="list-group-flush">
			<ListGroupItem className="border-bottom d-flex justify-content-between align-items-center bg-light cursor-pointer" onClick={TODO}>
				<div>New Chat</div>
				<FontAwesomeIcon icon={faPlus} className="me-2" />
			</ListGroupItem>
			{Object.values(chats).map(chat => {
				let bot = bots[chat.bot_id];
				return (
					<ListGroupItem 
						key={chat.ID} 
						className="border-bottom" 
						style={{ "cursor": "pointer" }}
						onClick={() => dispatch(goToBotChat(chat))}>
						<Container className="text-start">
							<div className="d-flex justify-content-between">
								<div className="flex-grow-1">
									<div><strong>{bot.name}</strong></div>
									<div>{chat.ID}</div>
								</div>
								<div className="d-flex align-items-center">
									<FontAwesomeIcon 
										icon={faChevronRight} 
										className="ms-2 cursor-pointer" 
										style={{"cursor": "pointer"}} 
										/>
								</div>
							</div>
						</Container>
					</ListGroupItem>
				);
			})}
		</ListGroup>
	);
}