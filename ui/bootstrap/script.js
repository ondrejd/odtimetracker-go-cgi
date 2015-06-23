// Copyright 2015 Ondrej Donek. All rights reserved.
// See LICENSE file for more details about licensing.

/**
 * Holds projects (something like simple cache).
 * @var {Object}
 */
var gProjects = {};

/**
 * Generate random string.
 *
 * @param {Integer} length
 * @returns {String}
 */
function getRandomString(length) {
	var text = "";
	var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

	for (var i=0; i<length; i++) {
		text += possible.charAt(Math.floor(Math.random() * possible.length));
	}

	return text;
}

/**
 * Capitalize the first letter of given string.
 *
 * @param {String} str
 * @returns {String}
 * @see http://stackoverflow.com/questions/1026069/capitalize-the-first-letter-of-string-in-javascript#answer-1026087
 */
function capitalizeFirstLetter(str) {
	return str.charAt(0).toUpperCase() + str.slice(1);
} // end capitalizeFirstLetter(str)

/**
 * Add message to the list.
 *
 * @param {String} type Message type ('danger','info','success','warning').
 * @param {String} msg  Text of the message.
 * @returns {void}
 *
 * @todo Method for updating messages count is called BEFORE the message is actually removed...
 * @todo Add icons according to message type (e.g. for type 'error' use 'glyphicon-exclamation-sign' etc.)
 */
function addMessage(type, msg) {
	$("#messagesCont").append(
		'<div class="alert alert-warning alert-dismissible" role="alert">' +
			'<button type="button" class="close" data-dismiss="alert" aria-label="Close">' +
				'<span aria-hidden="true">&times;</span>' +
			'</button>' +
			'<strong>' + capitalizeFirstLetter(type) + '!</strong> ' + msg +
		'</div>'
	);

	// Event handler for closing messages
	$('#messagesCont .alert').on('closed.bs.alert', function () {
		console.log("On alert close...")
	});

	updateMessagesCountBadge();
} // end addMessage(aType, aMessage)

/**
 * Update badge in message tab button with current number of the messages.
 *
 * @returns {void}
 */
function updateMessagesCountBadge() {
	$("#messagesCountBadge").html($("#messagesCont > div.alert").length);
} // end updateMessagesCountBadge()

/**
 * Display a simple form for stopping the activity.
 *
 * @param {Object} activity
 * @returns {void}
 */
function showActivityStopForm(activity) {
	$('#runningActivityName').html(activity.Name);
	$('#runningActivityDescription').html(activity.Description);
	$('#stopActivityForm input[name="aid"]').val(activity.ActivityId);
	$('#startActivityJumbotron').hide();
	$('#stopActivityJumbotron').show();
} // end showActivityStopForm(activity)

/**
 * Display a simple form for starting the activity.
 *
 * @returns {void}
 */
function showActivityStartForm() {
	$('#projectIdInput').val('');
	$('#activityInput').val('');
	$('#projectInput').val('');
	$('#tagsInput').val('');
	$('#descriptionInput').val('');
	$('#startActivityJumbotron').show();
	$('#stopActivityJumbotron').hide();
}

/**
 * Check if there is a running activity and display the result on the hometab.
 *
 * @returns {void}
 */
function checkRunningActivity() {
	$.getJSON('/GetRunningActivity', function (data) {
		if (data.ActivityId > 0) {
			// There is a running activity...
			showActivityStopForm(data);
		} else {
			// There is no running activity...
			showActivityStartForm();
		}
	});
} // end checkRunningActivity()

/**
 * Load existing activities.
 *
 * @returns {void}
 */
function loadActivities() {
	$.getJSON('/ListActivities', function(data) {
		var items = [];
		$.each(data, function(key, v) {
			var project = gProjects[v.ProjectId];
			var projectName = project.Name;
			var cls = (v.Stopped == '') ? ' class="active"' : '';
			var desc = ((v.Description != '') ? '<br><small>' + v.Description + '</small>' : '');
			items.push(
			'<tr' + cls + '>' +
				'<td class="col-id">' + v.ActivityId + "</td>" +
				'<td class="col-star"><span class="glyphicon glyphicon-star-empty" title="..."></span></td>' +
				'<td class="col-text"><h4>' + v.Name + desc + "</h4></td>" +
				'<td class="col-project">' + projectName + "</td>" +
				'<td class="col-datetime">' + v.Started + "</td>" +
				'<td class="col-datetime">' + v.Stopped + "</td>" +
				'<td class="col-control">' +
					'<div class="btn-group" role="group" aria-label="...">' +
						'<button class="btn btn-danger btn-xs" data-toggle="modal" ' +
						        'data-target="#RemoveActivityDlg" ' +
						        'data-projectId="' + v.ProjectId + '" ' +
						        'data-projectName="' + projectName + '" ' +
						        'data-name="' + v.Name + '" ' +
						        'data-description="' + v.Description.replace('"', '\"') + '" ' +
						        'data-tags="' + v.Tags + '" ' +
						        'data-started="' + v.Started + '" ' +
						        'data-stopped="' + v.Stopped + '" ' +
						        'title="Remove activity">' +
							'<span class="glyphicon glyphicon-remove-circle"></span>' +
						'</button>' +
						'<button class="btn btn-primary btn-xs" data-toggle="modal" ' +
						        'data-target="#EditActivityDlg" ' +
						        'data-projectId="' + v.ProjectId + '" ' +
						        'data-projectName="' + projectName + '" ' +
						        'data-name="' + v.Name + '" ' +
						        'data-description="' + v.Description.replace('"', '\"') + '" ' +
						        'data-tags="' + v.Tags + '" ' +
						        'data-started="' + v.Started + '" ' +
						        'data-stopped="' + v.Stopped + '" ' +
						        'title="Edit activity">' +
							'<span class="glyphicon glyphicon-pencil"></span>' +
						'</button>' +
						'<button class="btn btn-info btn-xs repeat-activity-dlg" ' +
						        'data-toggle="modal" ' +
						        'data-projectId="' + v.ProjectId + '" ' +
						        'data-projectName="' + projectName + '" ' +
						        'data-name="' + v.Name + '" ' +
						        'data-description="' + v.Description.replace('"', '\"') + '" ' +
						        'data-tags="' + v.Tags + '" ' +
						        'data-started="' + v.Started + '" ' +
						        'data-stopped="' + v.Stopped + '" ' +
						        'title="Repeat activity" ' +
						        'onclick="repeatActivity(event);">' +
							'<span class="glyphicon glyphicon-repeat"></span>' +
						'</button>' +
					'</div>' +
				'</td>' +
			'</tr>');
		});
		$('#activitiesTable').html(items.join(''));
	});
} // end loadActivities()

/**
 * Load existing activities.
 *
 * @returns {void}
 */
function loadProjects() {
	// Every-time we load projects re-create the cache:
	gProjects = {};
	$.getJSON('/ListProjects', function(data) {
		var items = [];
		$.each(data, function(key, v) {
			gProjects[v.ProjectId] = v;
			var desc = ((v.Description != '') ? '<br><small>' + v.Description + '</small>' : '');
			items.push(
			'<tr>' +
				'<td class="col-id">' + v.ProjectId + '</td>' +
				'<td class="col-star">&nbsp;</td>' +
				'<td class="col-text">' +
					'<h4>' + v.Name + desc + '</h4>' +
				'</td>' +
				'<td class="col-datetime">' + v.Created + '</td>' +
				'<td class="col-control">' +
					'<div class="btn-group" role="group" aria-label="...">' +
						'<button class="btn btn-danger btn-xs" data-toggle="modal" ' +
						        'data-target="#RemoveProjectDlg" ' +
						        'data-projectId="' + v.ProjectId + '" ' +
						        'title="Remove project">' +
							'<span class="glyphicon glyphicon-remove-circle"></span>' +
						'</button> ' +
						'<button class="btn btn-primary btn-xs" data-toggle="modal" ' +
						        'data-target="#EditProjectDlg" ' +
						        'data-projectId="' + v.ProjectId + '" ' +
						        'data-name="' + v.Name + '" ' +
						        'data-description="' + v.Description + '" ' +
						        'data-created="' + v.Created + '" ' +
						        'title="Edit project">' +
							'<span class="glyphicon glyphicon-pencil"></span>' +
						'</button>' +
					'</div>' +
				'</td>' +
			'</tr>');
		});
		$('#projectsTable').html(items.join(''));
	});
} // end loadProjects()

/**
 * Repeat activity.
 *
 * @param {DOMEvent} event
 * @returns {void}
 */
function repeatActivity(event) {
	var button = $(event.target) // Button that triggered the modal
	var projectId = button.attr('data-projectId');
	var projectName = button.attr('data-projectName');
	var name = button.attr('data-name');
	var desc = button.attr('data-description');
	var tags = button.attr('data-tags');

	//$('#projectIdInput').val(name);
	$('#activityInput').val(name);
	$('#projectInput').val(projectName);
	$('#tagsInput').val(tags);
	$('#descriptionInput').val(desc);
	$('#startActivityJumbotron').show();
	$('#stopActivityJumbotron').hide();

	$('#tabBtn1 a').tab('show');
} // end repeatActivity(event)

/**
 * Initialize our application.
 */
$(document).ready(function () {
	// Firstly load all projects
	$.getJSON('/ListProjects', function(data) {
		// Re-format given data array for quicker access...
		$.each(data, function(key, v) {
			gProjects[v.ProjectId] = v;
		});
	});

	// Home tab
	$('#tabBtn1 a').click(function (event) {
		event.preventDefault();
		checkRunningActivity();
		$(this).tab('show');
	});

	// Activities tab
	$('#tabBtn2 a').click(function (event) {
		event.preventDefault();
		loadActivities();
		$(this).tab('show');
	});

	// Organize projects/tags tab
	$('#tabBtn3 a').click(function (event) {
		event.preventDefault();
		loadProjects();
		$(this).tab('show');
	});

	// Project name auto-complete
	$('#projectInput').autocomplete({
		minLength: 0,
		source: function(request, response) {
			$.ajax({
				type: 'POST',
				url: '/ProjectNameAutocomplete',
				dataType: 'json',
				data: { term: request.term },
				success: function(data) {
					response($.map(data, function(item) {
						return { label: item.Name, value: item.Name };
					}));
				}
			});
		}
	});

	// Start activity form
	$("#startActivityForm").submit(function (event) {
		event.preventDefault();
		$.post('/StartActivity', $(this).serialize(), function (response) {
			// TODO We should also check response/request ID!!!
			console.log(response);
			if ("Error" in response) {
				addMessage("error", response.Error.Message);
			} else {
				addMessage("success", response.Result.Message);
			}
			checkRunningActivity();
		});
	});
	$("#activityInput").focus();

	// Stop activity form
	$('#stopActivityForm').submit(function (event) {
		event.preventDefault();
		$.post('/StopActivity', $(this).serialize(), function (response) {
			// TODO We should also check response/request ID!!!
			console.log(response);
			if ('Error' in response) {
				addMessage('error', response.Error.Message);
			} else {
				addMessage('success', response.Result.Message);
			}
			checkRunningActivity();
		});
	});

	// Edit project dialog
	$('#EditProjectDlg').on('show.bs.modal', function (event) {
		console.log("XXX #RemoveProjectDlg");
		var button = $(event.relatedTarget) // Button that triggered the modal
		var projectId = button.attr('data-projectId');
		var name = button.attr('data-name');
		var desc = button.attr('data-description');
		var created = button.attr('data-created');

		console.log(projectId);
		console.log(name);
		console.log(desc);
		console.log(created);

		var modal = $(this);
		modal.find('#epd_projectId').val(projectId);
		modal.find('#epd_name').val(name);
		modal.find('#epd_description').val(desc);
		modal.find('#epd_created').val(created);
	});

	// Edit project dialog
	$('#RemoveProjectDlg').on('show.bs.modal', function (event) {
		console.log("XXX #RemoveProjectDlg");
		var button = $(event.relatedTarget) // Button that triggered the modal
		var projectId = button.attr('data-projectId');

		console.log(projectId);

		// ...
	});

	// Check if there is a running activity
	checkRunningActivity();
});
