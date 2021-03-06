const GRANT_USER = 1;
const GRANT_PRODUCT_ADD = 2;
const GRANT_PRODUCT_EDIT = 4;
const GRANT_ITEM_READ = 8;
const GRANT_ITEM_ADD = 16;
const GRANT_ITEM_EDIT = 32;
const GRANT_DEL = 64;
const GRANT_SUPER = 127;


var $ = layui.$;

function init() {
	if (UINFO.Grant & GRANT_USER) {
		$('.WITH_GRANT_USER').removeClass("layui-hide")
	}
	if (UINFO.Grant & GRANT_PRODUCT_ADD) {
		$('.WITH_GRANT_PRODUCT_ADD').removeClass("layui-hide")
	}
	if (UINFO.Grant & GRANT_PRODUCT_EDIT) {
		$('.WITH_GRANT_PRODUCT_EDIT').removeClass("layui-hide")
	}
	if (UINFO.Grant & GRANT_ITEM_READ) {
		$('.WITH_GRANT_ITEM_READ').removeClass("layui-hide")
	}
	if (UINFO.Grant & GRANT_ITEM_ADD) {
		$('.WITH_GRANT_ITEM_ADD').removeClass("layui-hide")
	}
	if (UINFO.Grant & GRANT_ITEM_EDIT) {
		$('.WITH_GRANT_ITEM_EDIT').removeClass("layui-hide")
	}
	if (UINFO.Grant & GRANT_DEL) {
		$('.WITH_GRANT_DEL').removeClass("layui-hide")
	}
	if (UINFO.Grant & GRANT_SUPER) {
		$('.WITH_GRANT_SUPER').removeClass("layui-hide")
	}
}


layui.use(['element', 'form', 'layer', 'table', 'laytpl'], function () {
	var element = layui.element,
		form = layui.form,
		layer = layui.layer,
		table = layui.table;




	// 更改密码
	form.on('submit(change-password)', function (data) {
		if (data.field["new-passwd"] != data.field["new-passwd2"]) {
			layer.msg("新密码不一致")
			return false
		}
		var upload_data = {
			"old-password": data.field["old-passwd"],
			"new-password": data.field["new-passwd"]
		};
		xmlhttp.onreadystatechange = function () {
			if (xmlhttp.readyState == 4) {
				if (xmlhttp.status == 200) {
					var res = JSON.parse(xmlhttp.response)
					if (res.res) {
						layer.msg(res.msg + "1 秒后跳转", {
							time: 1000 //1秒关闭（如果不配置，默认是3秒）
						}, function () {
							location.href = "/backend/";
						})
					} else {
						layer.msg(res.msg);
					}
				} else {
					layer.msg("服务器连接失败：" + xmlhttp.status)
				}
			}
		}
		xmlhttp.open("POST", "/api/changepw/", true);
		xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
		xmlhttp.send(toURL(upload_data))
		return false
	})

	//用户管理
	table.render({
		elem: '#tb-user',
		height: 'full-300',
		url: '/api/select_user',
		cols: [[
			{ field: 'id', title: 'ID', width: '10%', sort: true },
			{ field: 'name', title: '用户名', width: '10%' },
			{
				field: 'grant', title: '权限', width: '60%', templet: function (d) {
					s = '';
					if (d.grant & GRANT_USER)
						s += '用户编辑 ';
					if (d.grant & GRANT_PRODUCT_ADD)
						s += '添加品类 ';
					if (d.grant & GRANT_PRODUCT_EDIT)
						s += '编辑品类 ';
					if (d.grant & GRANT_ITEM_READ)
						s += '读取商品 ';
					if (d.grant & GRANT_ITEM_ADD)
						s += '添加商品 ';
					if (d.grant & GRANT_ITEM_EDIT)
						s += '编辑商品 ';
					if (d.grant & GRANT_DEL)
						s += '删除条目 ';
					return s;
				}
			},
			{ field: 'option', title: '', width: '15%', toolbar: '#user-tool' },
		]],
		toolbar: '#user-toolbar',
		//	page: true,
	})
	$('#user-name-search').on("input", function (e) {
		setTimeout(function () {
			if ($('#user-name-search').val() != "") {
				table.reload('tb-user', {
					url: '/api/select_user/',
					where: {
						"search_name": '%' + $('#user-name-search').val() + '%',
					},
				}, true)
			} else {
				table.reload('tb-user', {
					url: '/api/select_user/',
					where: null,
				}, true)
			}
		}, 0)
	})
	table.on('toolbar(tb-user)', function (obj) {
		switch (obj.event) {
			case 'add':
				$("#user-layer-add-username").val("");
				layer.open({
					type: 1,
					content: $("#user-layer-add"),
					title: '添加用户',
					btn: '添加用户',
					resize: false,
					scrollbar: false,
					yes: function (index, layero) {
						var grant = 0;
						if ($("#user-layer-add-grant-user").prop("checked"))
							grant += GRANT_USER;
						if ($("#user-layer-add-grant-product-add").prop("checked")) {
							grant += GRANT_PRODUCT_ADD;
						}
						if ($("#user-layer-add-grant-product-edit").prop("checked")) {
							grant += GRANT_PRODUCT_EDIT;
						}
						if ($("#user-layer-add-grant-item-read").prop("checked")) {
							grant += GRANT_ITEM_READ;
						}
						if ($("#user-layer-add-grant-item-add").prop("checked")) {
							grant += GRANT_ITEM_ADD;
						}
						if ($("#user-layer-add-grant-item-edit").prop("checked")) {
							grant += GRANT_ITEM_EDIT;
						}
						if ($("#user-layer-add-grant-del").prop("checked")) {
							grant += GRANT_DEL;
						}
						var upload_data = {
							"name": $("#user-layer-add-username").val(),
							"password": $("#user-layer-add-password").val(),
							"grant": grant,
						};
						if (upload_data.name == "" || upload_data.password == "") {
							layer.msg("用户名或密码不能为空")
						} else {
							xmlhttp.onreadystatechange = function () {
								if (xmlhttp.readyState == 4) {
									if (xmlhttp.status == 200) {
										var res = JSON.parse(xmlhttp.response)
										layer.msg(res.msg);
										if (res.res) {
											$('#user-name-search').val("")
											table.reload('tb-user', {
												url: '/api/select_user/',
											}, true)
										}
									} else {
										layer.msg("服务器连接失败：" + xmlhttp.status)
									}
								}
							}
							xmlhttp.open("POST", "/api/add_user/", true);
							xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
							xmlhttp.send(toURL(upload_data))
							layer.close(index);
						}
					},
				});
				break;
		};
	});
	table.on('tool(tb-user)', function (obj) {
		switch (obj.event) {
			case 'edit':
				id = obj.data.id;
				grant = obj.data.grant;
				$("#user-layer-edit-grant-user").prop("checked", Boolean(grant & GRANT_USER));
				$("#user-layer-edit-grant-product-add").prop("checked", Boolean(grant & GRANT_PRODUCT_ADD));
				$("#user-layer-edit-grant-product-edit").prop("checked", Boolean(grant & GRANT_PRODUCT_EDIT));
				$("#user-layer-edit-grant-item-read").prop("checked", Boolean(grant & GRANT_ITEM_READ));
				$("#user-layer-edit-grant-item-add").prop("checked", Boolean(grant & GRANT_ITEM_ADD));
				$("#user-layer-edit-grant-item-edit").prop("checked", Boolean(grant & GRANT_ITEM_EDIT));
				$("#user-layer-edit-grant-del").prop("checked", Boolean(grant & GRANT_DEL));
				layui.form.render();
				layer.open({
					type: 1,
					content: $("#user-layer-edit"),
					title: '编辑用户',
					btn: '保存',
					resize: false,
					scrollbar: false,
					yes: function (index, layero) {
						var grant = 0;
						if ($("#user-layer-edit-grant-user").prop("checked"))
							grant += GRANT_USER;
						if ($("#user-layer-edit-grant-product-add").prop("checked")) {
							grant += GRANT_PRODUCT_ADD;
						}
						if ($("#user-layer-edit-grant-product-edit").prop("checked")) {
							grant += GRANT_PRODUCT_EDIT;
						}
						if ($("#user-layer-edit-grant-item-read").prop("checked")) {
							grant += GRANT_ITEM_READ;
						}
						if ($("#user-layer-edit-grant-item-add").prop("checked")) {
							grant += GRANT_ITEM_ADD;
						}
						if ($("#user-layer-edit-grant-item-edit").prop("checked")) {
							grant += GRANT_ITEM_EDIT;
						}
						if ($("#user-layer-edit-grant-del").prop("checked")) {
							grant += GRANT_DEL;
						}
						var upload_data = {
							"uid": id,
							"password": $("#user-layer-edit-password").val(),
							"grant": grant,
						};
						xmlhttp.onreadystatechange = function () {
							if (xmlhttp.readyState == 4) {
								if (xmlhttp.status == 200) {
									var res = JSON.parse(xmlhttp.response)
									layer.msg(res.msg);
									if (res.res) {
										$('#user-name-search').val("")
										table.reload('tb-user', {
											url: '/api/select_user/',
										}, true)
									}
								} else {
									layer.msg("服务器连接失败：" + xmlhttp.status)
								}
							}
						}
						xmlhttp.open("POST", "/api/set_user/", true);
						xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
						xmlhttp.send(toURL(upload_data))
						layer.close(index);
					},
				});
				break;
			case 'delete':
				layer.confirm('确定删除用户“' + obj.data.name + '”?', {
					btn: ['确定', '取消'] //按钮
				}, function () {
					var upload_data = {
						"uid": obj.data.id,
					};
					xmlhttp.onreadystatechange = function () {
						if (xmlhttp.readyState == 4) {
							if (xmlhttp.status == 200) {
								var res = JSON.parse(xmlhttp.response)
								layer.msg(res.msg);
								if (res.res) {
									$('#user-name-search').val("")
									table.reload('tb-user', {
										url: '/api/select_user/',
									}, true)
								}
							} else {
								layer.msg("服务器连接失败：" + xmlhttp.status)
							}
						}
					}
					xmlhttp.open("POST", "/api/del_user/", true);
					xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
					xmlhttp.send(toURL(upload_data))
				});
				break;
		};
	});
	table.on('sort(tb-user)', function (obj) { init() })

	//品类管理

	table.render({
		elem: '#tb-pattern',
		height: 'full-300',
		url: '/api/select_pattern',
		cols: [[
			{ field: 'id', title: 'ID', width: '10%', sort: true },
			{ field: 'name', title: '品名', width: '35%' },
			{ field: 'brand', title: '品牌', width: '15%' },
			{
				field: '', title: '款式', width: '10%', templet: function (d) {
					return '<button class="layui-btn layui-btn-xs" onclick="showProduct(' + d.id + ');">查看款式</button>'
				}
			},
			{ field: 'price', title: '价格', width: '10%', sort: true },
			{ field: 'option', title: '', width: '15%', toolbar: '#pattern-tool' },
		]],
		toolbar: '#pattern-toolbar',
		done: (res, curr, count) => { init(); },
		//	page: true,
	})
	$('#pattern-search').on("input", function (e) {
		setTimeout(function () {
			if ($('#pattern-search').val() != "") {
				table.reload('tb-pattern', {
					url: '/api/select_pattern/',
					where: {
						"search_key": '%' + $('#pattern-search').val() + '%',
					},
				}, true)
			} else {
				table.reload('tb-pattern', {
					url: '/api/select_pattern/',
					where: null,
				}, true)
			}
		}, 0)
	})

	table.on('toolbar(tb-pattern)', function (obj) {
		switch (obj.event) {
			case 'add':
				layer.open({
					type: 1,
					content: $("#pattern-layer-add"),
					title: '添加品类',
					btn: '添加品类',
					resize: false,
					scrollbar: false,
					yes: function (index, layero) {
						var upload_data = {
							"name": $("#pattern-layer-add-name").val(),
							"brand": $("#pattern-layer-add-brand").val(),
							"price": $("#pattern-layer-add-price").val(),
						};
						if (upload_data.name == "" || upload_data.brand == "" || upload_data.price == "") {
							layer.msg("品名或品牌或价格不能为空")
						} else {
							xmlhttp.onreadystatechange = function () {
								if (xmlhttp.readyState == 4) {
									if (xmlhttp.status == 200) {
										var res = JSON.parse(xmlhttp.response)
										layer.msg(res.msg);
										if (res.res) {
											$('#pattern-search').val("")
											table.reload('tb-pattern', {
												url: '/api/select_pattern/',
											}, true)
										}
									} else {
										layer.msg("服务器连接失败：" + xmlhttp.status)
									}
								}
							}
							xmlhttp.open("POST", "/api/add_pattern/", true);
							xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
							xmlhttp.send(toURL(upload_data))
							layer.close(index);
						}
					},
				});
				break;
		};
	});

	table.on('tool(tb-pattern)', function (obj) {
		switch (obj.event) {
			case 'add':
				layer.open({
					type: 1,
					content: $("#pattern-layer-add-product"),
					title: '添加款式',
					btn: '添加款式',
					resize: false,
					scrollbar: false,
					yes: function (index, layero) {
						var upload_data = {
							"id": obj.data.id,
							"SKU": $("#pattern-layer-add-SKU").val(),
							"color": $("#pattern-layer-add-color").val(),
							"size": $("#pattern-layer-add-size").val(),
						};
						if (false) {
						} else {
							xmlhttp.onreadystatechange = function () {
								if (xmlhttp.readyState == 4) {
									if (xmlhttp.status == 200) {
										var res = JSON.parse(xmlhttp.response)
										layer.msg(res.msg);

									} else {
										layer.msg("服务器连接失败：" + xmlhttp.status)
									}
								}
							}
							xmlhttp.open("POST", "/api/add_product/", true);
							xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
							xmlhttp.send(toURL(upload_data))
							layer.close(index);
						}
					},
				});
				break;
			case 'delete':
				layer.confirm('确定删除品类“' + obj.data.name + '”?', {
					btn: ['确定', '取消'] //按钮
				}, function () {
					var upload_data = {
						"id": obj.data.id,
					};
					xmlhttp.onreadystatechange = function () {
						if (xmlhttp.readyState == 4) {
							if (xmlhttp.status == 200) {
								var res = JSON.parse(xmlhttp.response)
								layer.msg(res.msg);
								if (res.res) {
									$('#pattern-search').val("")
									table.reload('tb-pattern', {
										url: '/api/select_pattern/',
									}, true)
								}
							} else {
								layer.msg("服务器连接失败：" + xmlhttp.status)
							}
						}
					}
					xmlhttp.open("POST", "/api/del_pattern/", true);
					xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
					xmlhttp.send(toURL(upload_data))
				});
				break;
		};
	});
	table.on('sort(tb-pattern)', function (obj) { init() })
	// 商品管理
	table.render({
		elem: '#tb-batch',
		height: 'full-300',
		url: '/api/select_batch',
		cols: [[
			{ field: 'id', title: 'ID', width: '10%', sort: true },
			{ field: 'uid', title: '管理员ID', width: '10%', sort: true },
			{ field: 'time', title: '时间', width: '30%', sort: true },
			{ field: 'count', title: '数量', width: '15%', sort: true },
			{ field: 'option', title: '', width: '30%', toolbar: '#batch-tool' },
		]],
		toolbar: '#batch-toolbar',
		done: (res, curr, count) => { init(); },
		//	page: true,
	})

	table.on('toolbar(tb-batch)', function (obj) {
		switch (obj.event) {
			case 'add':
				$('#batch-layer-add-pattern').html("<option value = '' >请选择</option>");
				$('#batch-layer-add-product').html("");
				batch_layer_add_upload = {}
				updateBatchUpload()
				form.render(null, 'form-batch-layer-add');
				xmlhttp.onreadystatechange = function () {
					if (xmlhttp.readyState == 4) {
						if (xmlhttp.status == 200) {
							var res = JSON.parse(xmlhttp.response)
							if (res.code == 0) {
								var h = "<option value = '' >请选择</option>";
								for (i in res.data) {
									h += "<option value = '" + res.data[i].id + "' >" + res.data[i].brand + " - " + res.data[i].name + "</option>";
								}
								$('#batch-layer-add-pattern').html(h);
								form.render(null, 'form-batch-layer-add');
							} else {
								layer.msg(res.msg)
							}
						} else {
							layer.msg("服务器连接失败：" + xmlhttp.status)
						}
					}
				}
				xmlhttp.open("GET", "/api/select_pattern/", true);
				// xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
				xmlhttp.send()

				layer.open({
					type: 1,
					content: $("#batch-layer-add"),
					title: '添加商品',
					btn: '添加商品',
					// resize: false,
					// scrollbar: false,
					yes: function (index, layero) {
						for (i in batch_layer_add_upload) {
							num = batch_layer_add_upload[i]["num"]
							batch_layer_add_upload[i] = num
						}
						var upload_data = {
							data: JSON.stringify(batch_layer_add_upload)
						};
						if (Object.keys(batch_layer_add_upload).length == 0) {
							layer.msg("商品不能为空")
						} else {
							xmlhttp.onreadystatechange = function () {
								if (xmlhttp.readyState == 4) {
									if (xmlhttp.status == 200) {
										var res = JSON.parse(xmlhttp.response)
										layer.msg(res.msg);
										if (res.res) {
											$('#pattern-search').val("")
										}
										table.reload('tb-batch', {
											url: '/api/select_batch/',
										}, true)

									} else {
										layer.msg("服务器连接失败：" + xmlhttp.status)
									}
								}
							}
							xmlhttp.open("POST", "/api/add_batch/", true);
							xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
							xmlhttp.send(toURL(upload_data))
							layer.close(index);
						}
					},
				});
				break;
		};
	});
	table.render({
		elem: '#tb-batch-info',
		height: 'full-300',
		url: '/api/info_batch',
		cols: [[
			{ field: 'pd_id', title: 'ID', width: '100', sort: true },
			{ field: 'pt_brand', title: '品牌', width: '100' },
			{ field: 'pt_name', title: '品名', width: '300' },
			{ field: 'pd_color', title: '颜色', width: '300' },
			{ field: 'pd_size', title: '尺码', width: '300' },
			{ field: 'count', title: '数量', width: '150', sort: true },
			// { field: 'option', title: '', width: '30%', toolbar: '#batch-tool' },
		]],
		done: (res, curr, count) => { init(); },
		//	page: true,
	})
	table.on('tool(tb-batch)', function (obj) {
		switch (obj.event) {
			case 'info':
				table.reload("tb-batch-info", {
					url: '/api/info_batch/',
					where: {
						"id": obj.data.id,
					},
				}, true);
				layer.open({
					type: 1,
					content: $("#batch-layer-info"),
					title: '批次详情',
					// resize: false,
					// scrollbar: false,
				});
				break;
			case 'delete':
				layer.confirm('确定删除批次“' + obj.data.id + '”?', {
					btn: ['确定', '取消'] //按钮
				}, function () {
					var upload_data = {
						"id": obj.data.id,
					};
					xmlhttp.onreadystatechange = function () {
						if (xmlhttp.readyState == 4) {
							if (xmlhttp.status == 200) {
								var res = JSON.parse(xmlhttp.response)
								layer.msg(res.msg);
								if (res.res) {
									table.reload('tb-batch', {
										url: '/api/select_batch/',
									}, true)
								}
							} else {
								layer.msg("服务器连接失败：" + xmlhttp.status)
							}
						}
					}
					xmlhttp.open("POST", "/api/del_batch/", true);
					xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
					xmlhttp.send(toURL(upload_data))
				});
				break;
		};
	});
	form.on('select(batch-layer-add-pattern)', function (data) {
		if (data.value == "") {
			$('#batch-layer-add-product').html("");
			form.render(null, 'form-batch-layer-add');
			return
		}
		var upload_data = {
			"ptid": data.value
		}
		xmlhttp.onreadystatechange = function () {
			if (xmlhttp.readyState == 4) {
				if (xmlhttp.status == 200) {
					var res = JSON.parse(xmlhttp.response)
					if (res.code == 0) {
						var h = "<option value = '' >请选择</option>";
						for (i in res.data) {
							h += "<option value = '" + res.data[i].id + "' >" + res.data[i].SKU + " - " + res.data[i].color + " - " + res.data[i].size + "</option>";
						}
						$('#batch-layer-add-product').html(h);
						form.render(null, 'form-batch-layer-add');
					} else {
						layer.msg(res.msg)
					}
				} else {
					layer.msg("服务器连接失败：" + xmlhttp.status)
				}
			}
		}
		xmlhttp.open("GET", "/api/select_product/?" + toURL(upload_data), true);
		// xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
		xmlhttp.send()
	});
	form.on("submit(batch-layer-add-add)", (data) => {
		id = $("#batch-layer-add-product").val();
		num = Number($("#batch-layer-add-count").val());

		if (id == null || id == "" || num == 0)
			return

		// console.log(id, num)
		text = data.form[1].value + " - " + data.form[3].value;
		if (batch_layer_add_upload[id]) {
			batch_layer_add_upload[id].num += num;
		} else {
			batch_layer_add_upload[id] = { "text": text, "num": num };
		}

		updateBatchUpload();
	})

	// item
	table.render({
		elem: '#tb-item',
		height: 'full-300',
		url: '/api/select_item',
		cols: [[
			{ field: 'bt_id', title: '批次号', width: '10%' },
			{ field: 'pt_name', title: '品名', width: '10%' },
			{ field: 'pt_brand', title: '品牌', width: '10%' },
			{ field: 'pd_SKU', title: 'SKU', width: '10%' },
			{ field: 'pd_color', title: '颜色', width: '10%' },
			{ field: 'pd_size', title: '规格', width: '10%' },
			{ field: 'pt_price', title: '价格', width: '10%' },
			{ field: 'it_id', title: '识别号', width: '25%' },
		]],
		toolbar: '#item-toolbar',
		page: true,
	})
	$("#item-search-batch").on("input", reloadItemTB)
});
function reloadItemTB() {
	batch = $("#item-search-batch").val()
	SKU = $("#item-search-SKU").val()
	id = $("#item-search-id").val()
	layui.table.reload('tb-item', {
		url: '/api/select_item/',
		where: {
			"id": id,
			"SKU": SKU,
			"batch": batch,
		},
	}, true)
}

function updateBatchUpload() {
	html = ""
	for (i in batch_layer_add_upload) {
		html += "<tr>"
		html += "<td>" + i + "</td>";
		html += "<td>" + batch_layer_add_upload[i]["text"] + "</td>";
		html += "<td>" + batch_layer_add_upload[i]["num"] + "</td>";
		html += "<td><button class='layui-btn layui-btn-sm' onclick='delete batch_layer_add_upload[" + i + "];updateBatchUpload();'>删除</button></td>";
		html += "</tr>"
	}
	$("#batch-layer-add-data").html(html);
	layui.table.init('batch-layer-add-tb')
}



function showContent(select) {
	$(".body-content").addClass("layui-hide");
	if (select != "")
		$(select).removeClass("layui-hide");
}

function showBatch(id) {
	var upload_data = {
		"id": id,
	};
	xmlhttp.onreadystatechange = function () {
		if (xmlhttp.readyState == 4) {
			if (xmlhttp.status == 200) {
				var res = JSON.parse(xmlhttp.response)
				// var tpl = $('#pattern-layer-show-productTpl').html();
				// layui.use('laytpl', function () {
				// 	var laytpl = layui.laytpl;
				// 	laytpl(tpl).render(res, function (html) {
				// 		$("#pattern-layer-show-product").html(html);
				// 		layui.table.init('tb-product', {
				// 			done: (res, curr, count) => { init(); }
				// 		});
				// 		layui.layer.open({
				// 			type: 1,
				// 			title: '款式列表',
				// 			content: $("#pattern-layer-show-product"),
				// 			resize: false,
				// 		});
				// 	});
				// })

			} else {
				layer.msg("服务器连接失败：" + xmlhttp.status)
			}
		}
	}
	// xmlhttp.open("GET", "/api/select_product/?" + toURL(upload_data), true);
	xmlhttp.send()
}


function showProduct(id) {
	var upload_data = {
		"ptid": id,
	};
	xmlhttp.onreadystatechange = function () {
		if (xmlhttp.readyState == 4) {
			if (xmlhttp.status == 200) {
				var res = JSON.parse(xmlhttp.response)
				var tpl = $('#pattern-layer-show-productTpl').html();
				layui.use('laytpl', function () {
					var laytpl = layui.laytpl;
					laytpl(tpl).render(res, function (html) {
						$("#pattern-layer-show-product").html(html);
						layui.table.init('tb-product', {
							done: (res, curr, count) => { init(); }
						});
						layui.layer.open({
							type: 1,
							title: '款式列表',
							content: $("#pattern-layer-show-product"),
							resize: false,
						});
					});
				})

			} else {
				layer.msg("服务器连接失败：" + xmlhttp.status)
			}
		}
	}
	xmlhttp.open("GET", "/api/select_product/?" + toURL(upload_data), true);
	//xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
	xmlhttp.send()
}
function delProduct(id) {
	layer.confirm('确定删除款式 ' + id + '?', {
		btn: ['确定', '取消'] //按钮
	}, function () {
		var upload_data = {
			"id": id,
		};
		xmlhttp.onreadystatechange = function () {
			if (xmlhttp.readyState == 4) {
				if (xmlhttp.status == 200) {
					var res = JSON.parse(xmlhttp.response)
					layer.msg(res.msg);

					layui.use('laytpl', function () {
						var laytpl = layui.laytpl;
						var tpl = $('#pattern-layer-show-productTpl').html();
						laytpl(tpl).render(res, function (html) {
							$("#pattern-layer-show-product").html(html);
							layui.table.init('tb-product', {
								done: (res, curr, count) => { init(); }
							});
						});
					})

				} else {
					layer.msg("服务器连接失败：" + xmlhttp.status)
				}
			}
		}
		xmlhttp.open("POST", "/api/del_product/", true);
		xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
		xmlhttp.send(toURL(upload_data))
	});
}