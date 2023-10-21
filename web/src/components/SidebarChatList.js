import { useDispatch, useSelector } from "react-redux";
import { Button, Container, ListGroup, ListGroupItem } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faChevronLeft, faChevronRight, faPlus } from "@fortawesome/free-solid-svg-icons";
import { goToBotChat, goToBotEdit, goToSidebarBotList, selectSelectedBot } from "../store/page";
import { addChat, selectChats } from "../store/chats";

export default function SidebarChatList() {
	const dispatch = useDispatch();
	const chats = useSelector(selectChats);
	const bot = useSelector(selectSelectedBot);
	
	return (
		<ListGroup className="list-group-flush bg-dark text-white">
			<ListGroupItem className="border-bottom bg-dark text-white">
				<Container className="text-center">
					<div className="d-flex align-items-center justify-content-between mb-3">
						<div>
							{bot.name}
						</div>
						<div>
							<Button className="btn-sm ms-2" onClick={() => dispatch(goToBotEdit(bot))}>
								Edit
							</Button>
						</div>
					</div>
					<div>
						<Button className="btn-sm me-2" onClick={() => {dispatch(goToSidebarBotList())}}>
							<FontAwesomeIcon icon={faChevronLeft} className="me-2" />
							Back to Bot List
						</Button>
						<Button className="btn-sm" onClick={() => dispatch(addChat(bot.ID))}>
							<FontAwesomeIcon icon={faPlus} className="me-2" />
							New Chat
						</Button>
					</div>
				</Container>
			</ListGroupItem>
			{chats.map(chat => {
				return (
					<ListGroupItem 
						key={chat.ID} 
						className="bg-dark text-white border-bottom" 
						style={{ "cursor": "pointer" }}
						onClick={() => dispatch(goToBotChat(chat))}>
						<Container className="text-start">
							<div className="d-flex justify-content-between">
								<div className="flex-grow-1">
									<strong>{chat.ID}</strong>
									<div>{chat.date_created}</div>
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