const GRANT_USER = 1;
const GRANT_PRODUCT_ADD = 2;
const GRANT_PRODUCT_EDIT = 4;
const GRANT_ITEM_READ = 8;
const GRANT_ITEM_ADD = 16;
const GRANT_ITEM_EDIT = 32;


var $ = layui.$;

layui.use(['element', 'form', 'layer', 'table'], function () {
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
			{ field: 'id', title: 'ID', width: '20%', sort: true },
			{ field: 'name', title: '用户名', width: '20%' },
			{ field: 'grant', title: '权限', width: '40%' },
			{ field: 'edit', width: '15%', toolbar: '#user-tool' },
		]],
		toolbar: '#user-toolbar',
		page: true,
	})
	$('#user-name-search').on("input", function (e) {
		setTimeout(function () {
			table.reload('tb-user', {
				url: '/api/select_user/',
				where: {
					"search_name": '%' + $('#user-name-search').val() + '%',
				},
			}, true)
		}, 0)
	})
	table.on('toolbar(tb-user)', function (obj) {
		switch (obj.event) {
			case 'add':
				$("#user-layer-username").val("");
				layer.open({
					type: 1,
					content: $("#user-layer-add"),
					title: '添加用户',
					btn: '添加用户',
					resize: false,
					scrollbar: false,
					yes: function (index, layero) {
						var grant = 0;
						if ($("#user-layer-grant-user").prop("checked"))
							grant += GRANT_USER;
						if ($("#user-layer-grant-product-add").prop("checked")) {
							grant += GRANT_PRODUCT_ADD;
						}
						if ($("#user-layer-grant-product-edit").prop("checked")) {
							grant += GRANT_PRODUCT_EDIT;
						}
						if ($("#user-layer-grant-item-read").prop("checked")) {
							grant += GRANT_ITEM_READ;
						}
						if ($("#user-layer-grant-item-add").prop("checked")) {
							grant += GRANT_ITEM_ADD;
						}
						if ($("#user-layer-grant-item-edit").prop("checked")) {
							grant += GRANT_ITEM_EDIT;
						}
						var upload_data = {
							"name": $("#user-layer-username").val(),
							"password": $("#user-layer-password").val(),
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
		console.log(obj)
		switch (obj.event) {
			case 'edit':
				layer.msg("setuser todo");
				break;
			case 'delete':
				layer.msg("deluser todo");
				break;
		};
	});

});
function showContent(select) {
	$(".body-content").addClass("layui-hide");
	if (select != "")
		$(select).removeClass("layui-hide");
}