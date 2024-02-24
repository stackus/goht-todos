htmx.onLoad(function (content) {
	const sortables = content.querySelectorAll(".sortable");
	for (let i = 0; i < sortables.length; i++) {
		const sortable = sortables[i];
		new Sortable(sortable, {
      draggable: '.draggable',
			handle: '.handle',
      animation: 150,
      chosenClass: 'dragClass'
    });
  }
});
