// Copyright 2015 Ondrej Donek. All rights reserved.
// See LICENSE file for more details about licensing.

/**
 * Capitalize the first letter of given string.
 *
 * @param {String} aStr
 * @returns {String}
 * @see http://stackoverflow.com/questions/1026069/capitalize-the-first-letter-of-string-in-javascript#answer-1026087
 */
function capitalizeFirstLetter(aStr) {
	return aStr.charAt(0).toUpperCase() + aStr.slice(1);
} // end capitalizeFirstLetter(aStr)

/**
 * Add message to the list.
 *
 * @param {String} aType    Message type ('danger','info','success','warning').
 * @param {String} aMessage Text of the message.
 * @returns {void}
 */
function addMessage(aType, aMessage) {
	$("#messagesCont").append(
		'<div class="alert alert-warning alert-dismissible" role="alert">' +
//			'<button type="button" class="close" data-dismiss="alert" aria-label="Close">' +
//				'<span aria-hidden="true" onclick="setTimeout(function() { updateMessagesCountBadge() }, 1200)">&times;</span>' +
//			'</button>' +
			'<strong>' + capitalizeFirstLetter(aType) + '!</strong> ' + aMessage +
		'</div>'
	);
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
 * @param {Object} aActivity
 * @returns {void}
 */
function printActivityStopForm(aActivity) {
	$("#runningActivityCont").html(
		'<div class="jumbotron">' +
			'<h1>' +
				aActivity.Name + '<br>' +
				'<small>' + aActivity.Description + '</small>' +
			'</h1>' +
			'<form id="stopActivityForm" class="form-horizontal">' +
				'<input type="hidden" name="aid" value="' + aActivity.ActivityId + '">' +
				'<div class="form-group col-sm-12">' +
					'<div class="col-sm-12">' +
						'<button type="submit" class="btn btn-primary btn-lg">Stop activity</button>' +
					'</div>' +
				'</div>' +
			'</form>' +
			'<p> <br> </p>' +
		'</div>'
	);
	$("#stopActivityForm").submit(function (event) {
		event.preventDefault();
		$.post('/StopActivity', $(this).serialize(), function (response) {
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
} // end printActivityStopForm(aActivity)

/**
 * Display for for inserting a new activity on the home tab.
 *
 * @returns {void}
 */
function printActivityForm() {
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
} // end printActivityForm()

/**
 * Check if there is a running activity and display the result on the hometab.
 *
 * @returns {void}
 */
function checkRunningActivity() {
	$.getJSON("/GetRunningActivity", function (data) {
		if (data.ActivityId > 0) {
			// There is a running activity...
			printActivityStopForm(data);
		} else {
			// There is no running activity...
			printActivityForm();
		}
	});
} // end checkRunningActivity()

/**
 * Load existing activities.
 *
 * @returns {void}
 */
function loadActivities() {
	$.getJSON("/ListActivities", function (data) {
		var items = [];
		$.each(data, function(key, v) {
			var cls = (v.Stopped == "") ? ' class="active"' : '';
			var desc = ((v.Description != "") ? "<br/><small>" + v.Description + "</small>" : "");
			items.push("<tr" + cls + ">" +
				'<td class="col-id">' + v.ActivityId + "</td>" +
				'<td class="col-star"><span class="glyphicon glyphicon-star-empty" title="..."></span></td>' +
				'<td class="col-text"><h4>' + v.Name + desc + "</h4></td>" +
				'<td class="col-project">' + v.ProjectId + "</td>" +
				'<td class="col-datetime">' + v.Started + "</td>" +
				'<td class="col-datetime">' + v.Stopped + "</td>" +
			"</tr>");
		});
		$("#activitiesTable").html(items.join(""));
	});
} // end loadActivities()

/**
 * Load existing activities.
 *
 * @returns {void}
 */
function loadProjects() {
	$.getJSON("/ListProjects", function (data) {
		var items = [];
		$.each(data, function(key, v) {
			var desc = ((v.Description != "") ? "<br/><small>" + v.Description + "</small>" : "");
			items.push("<tr>" +
				'<td class="col-id">' + v.ProjectId + "</td>" +
				'<td class="col-star">&nbsp;</td>' +
				'<td class="col-text"><h4>' + v.Name + desc + "</td>" +
				'<td class="col-datetime">' + v.Created + "</td>" +
			"</tr>");
		});
		$("#projectsTable").html(items.join(""));
	});
} // end loadProjects()

/**
 * Initialize our application.
 */
$(document).ready(function (e) {
	// Check if there is a running activity
	checkRunningActivity();

	// Home tab
	$('#tabBtn1 a').click(function (e) {
		e.preventDefault();
		checkRunningActivity();
		$(this).tab('show');
	});

	// Activities tab
	$('#tabBtn2 a').click(function (e) {
		e.preventDefault();
		loadActivities();
		$(this).tab('show');
	});

	// Organize projects/tags tab
	$('#tabBtn3 a').click(function (e) {
		e.preventDefault();
		loadProjects();
		$(this).tab('show');
	});

	// Helper tab with tasks
	$('#tabBtn4 a').click(function (e) {
		e.preventDefault();
		$(this).tab('show');
	});
});
