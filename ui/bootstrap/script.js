// Copyright 2015 Ondrej Donek. All rights reserved.
// See LICENSE file for more informations about licensing.

function checkRunningActivity() {
	$.getJSON("/GetRunningActivity", function (data) {
		if (!("ActivityId" in data)) {
			// Something goes wrong with the server...
			console.log(data);
			return;
		}
		if (data.ActivityId > 0) {
			console.log(data);
			// There is a running activity...
			$("#runningActivityCont").html(
				'<div class="jumbotron">' +
					'<h1>' + data.Name + '<br><small>' + data.Description + '</small></h1>' +
					'<form id="stopActivityForm" class="form-horizontal">' +
						'<input type="hidden" name="aid" value="' + data.ActivityId + '">' +
						'<div class="form-group col-sm-12">' +
							'<div class="col-sm-12">' +
								'<button type="submit" class="btn btn-primary btn-lg">Stop activity</button>' +
							'</div>' +
						'</div>' +
					'</form>' +
					'<p> <br> </p>' +
				'</div>'
			);
			// Stop activity form
			$("#stopActivityForm").submit(function (event) {
				console.log("Handler for [#stopActivityForm].submit() called.");
				var serializedData = $(this).serialize();
				console.log(serializedData);
				$.post('/StopActivity', serializedData, function (response) {
					console.log(response);
				});
				event.preventDefault();
			});
		} else {
			$("#runningActivityCont").html(
				'<div class="jumbotron">' +
					'<h1>No activity is running...</h1>' +
					'<p>You can start one right now!</p>' +
					'<form id="startActivityForm" class="form-horizontal">' +
						'<div class="form-group form-group-lg">' +
							'<label class="col-sm-2 control-label" for="activityInput">Activity</label>' +
							'<div class="col-sm-10">' +
								'<input class="form-control" type="text" name="name" id="activityInput" placeholder="Enter activity description&hellip;">' +
							'</div>' +
						'</div>' +
						'<div class="form-group form-group-sm">' +
							'<label class="col-sm-2 control-label" for="projectInput">Project</label>' +
							'<div class="col-sm-10">' +
								'<input class="form-control" type="text" name="project" id="projectInput" placeholder="Enter project&hellip;">' +
							'</div>' +
						'</div>' +
						'<div class="form-group form-group-sm">' +
							'<label class="col-sm-2 control-label" for="tagsInput">Tags</label>' +
							'<div class="col-sm-10">' +
								'<input class="form-control" type="text" name="tags" id="tagsInput" placeholder="Enter comma-separated tags&hellip;">' +
							'</div>' +
						'</div>' +
						'<div class="form-group form-group-sm">' +
							'<label class="col-sm-2 control-label" for="descInput">Description</label>' +
							'<div class="col-sm-10">' +
								'<input class="form-control" type="text" name="desc" id="descInput" placeholder="Enter additional description for the activity&hellip;">' +
							'</div>' +
						'</div>' +
						'<div class="form-group col-sm-12">' +
							'<div class="col-sm-2"> </div>' +
							'<div class="col-sm-10">' +
								'<button type="submit" class="btn btn-primary btn-lg">Start activity</button>' +
							'</div>' +
						'</div>' +
					'</form>' +
					'<p> <br> </p>' +
				'</div>'
			);
			// Set correct focus
			$("#activityInput").focus();
			// Start activity form
			$("#startActivityForm").submit(function (event) {
				var serializedData = $(this).serialize();
				$.post('/StartActivity', serializedData, function (response) {
					// TODO If is returned correct new activity show it as running
					// TODO Otherwise show error message
					console.log(response);
				});
				event.preventDefault();
			});
		}
	});
}

$(document).ready(function (e) {
	// Home tab
	$('#tabBtn1 a').click(function (e) {
		e.preventDefault();
		// Check if there is a running activity
		checkRunningActivity();
		// Show tab
		$(this).tab('show');
	});
	// Tab with activities list
	$('#tabBtn2 a').click(function (e) {
		e.preventDefault();
		// Load activities JSON
		$.getJSON("/ListActivities", function (data) {
			var items = [];
			$.each(data, function(key, v) {
				var cls = (v.Stopped == "") ? ' class="active"' : '';
				items.push("<tr" + cls + ">" +
					"<td>" + v.ActivityId + "</td>" +
					"<td>" + v.Name + "</td>" +
					"<td>" + v.Started + "</td>" +
					"<td>" + v.Stopped + "</td>" +
				"</tr>");
			});
			$("#activitiesTable").html(items.join(""));
		});
		// Show tab
		$(this).tab('show');
	});
	// Tab with projects and tags
	$('#tabBtn3 a').click(function (e) {
		e.preventDefault();
		// Show tab
		$(this).tab('show');
		// Load projects JSON
		$.getJSON("/ListProjects", function (data) {
			var items = [];
			$.each(data, function(key, v) {
				items.push("<tr>" +
					"<td>" + v.ProjectId + "</td>" +
					"<td>" + v.Name + "</td>" +
					"<td>" + v.Description + "</td>" +
					"<td>" + v.Created + "</td>" +
				"</tr>");
			});
			$("#projectsTable").html(items.join(""));
		});
	});
	// Helper tab with tasks
	$('#tabBtn4 a').click(function (e) {
		e.preventDefault();
		$(this).tab('show');
	});

	// Check if there is a running activity
	checkRunningActivity();
});
