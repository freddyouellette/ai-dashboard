import Swal from 'sweetalert2'

function confirm(text) {
	return new Promise((resolve) => 
		Swal.fire({
			title: text,
			confirmButtonText: "Confirm",
			reverseButtons: true,
			showCancelButton: true,
			heightAuto: false,
		}).then((result) => resolve(result.isConfirmed))
	);
}

function loading(text) {
	Swal.fire({
		title: text,
		allowOutsideClick: false,
		showConfirmButton: false,
		willOpen: () => {
			Swal.showLoading();
		},
	});
	return () => Swal.close();
}

const popup = { confirm, loading };

export default popup;