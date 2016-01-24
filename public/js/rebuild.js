$('.btn-rebuild').on('click', function () {
	$.ajax('/project/' + $(this).data('id') + '/build', {
		method: 'POST',
		success: function () {
			window.location.reload();
		},
		error: function () {
			alert('There was an error trying to rebuild your project.');
		}
	});
});
