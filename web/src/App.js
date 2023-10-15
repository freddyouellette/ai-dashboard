import { Button, Col, Container, ListGroup, ListGroupItem, Row } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faPlus } from '@fortawesome/free-solid-svg-icons'
import './App.css';

function App() {
	return (
		<Row className="App h-100">
			<Col className="sidebar text-white bg-dark h-100 p-0">
				<Container>
					<h1>AI Dashboard</h1>
				</Container>
				<ListGroup className="list-group-flush">
					<ListGroupItem className="bg-dark border-bottom">
						<Container className="text-end">
							<Button>
								<FontAwesomeIcon icon={faPlus} className="me-2" />
								Add New Bot
							</Button>
						</Container>
					</ListGroupItem>
					<ListGroupItem className="bg-dark text-white border-bottom">
						<Container className="text-start">
							Bot Name
						</Container>
					</ListGroupItem>
				</ListGroup>
			</Col>
			<Col className="col-8 h-100">
				<Container>
					Content
				</Container>
			</Col>
		</Row>
	);
}

export default App;
