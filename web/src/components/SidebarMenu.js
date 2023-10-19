import { useDispatch, useSelector } from "react-redux";
import { selectBots } from "../store/bots";
import { Button, Container, ListGroup, ListGroupItem } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faEdit, faPlus } from "@fortawesome/free-solid-svg-icons";
import { setBotToUpdate } from "../store/page";

export default function SidebarMenu() {
	const dispatch = useDispatch();
	const bots = useSelector(selectBots);
	
	return (
		<ListGroup className="list-group-flush">
			<ListGroupItem className="bg-dark border-bottom">
				<Container className="text-center">
					<Button onClick={() => {dispatch(setBotToUpdate(null))}}>
						<FontAwesomeIcon icon={faPlus} className="me-2" />
						Add New Bot
					</Button>
				</Container>
			</ListGroupItem>
			{bots.map(bot => {
				return (
					<ListGroupItem key={bot.ID} className="bg-dark text-white border-bottom">
						<Container className="text-start">
							<div className="d-flex justify-content-between align-items-center">
								<strong>{bot.name}</strong>
								<FontAwesomeIcon 
									icon={faEdit} 
									className="ms-2 cursor-pointer" 
									style={{"cursor": "pointer"}} 
									onClick={() => {dispatch(setBotToUpdate(bot))}} />
							</div>
							<div>{bot.description}</div>
						</Container>
					</ListGroupItem>
				);
			})}
		</ListGroup>
	);
}