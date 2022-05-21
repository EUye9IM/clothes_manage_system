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
		height: 'full-200',
		url: '/api/select_user',
		cols: [[
			{ field: 'id', title: 'ID', width: '20%', sort: true },
			{ field: 'name', title: '用户名', width: '20%' },
			{ field: 'grant', title: '权限', width: '40%' },
		]],
		toolbar: 'default',
		page: true,
	})
	$('#user-name-search').on("input", function (e) {
		setTimeout(function () {
			table.reload('tb-user', {
				url: '/api/select_user/',
				where: {
					"search_name": '%'+$('#user-name-search').val()+'%',
				},
			}, true)
		},0)
	})

});
function showContent(select) {
	$(".body-content").addClass("layui-hide");
	if (select != "")
		$(select).removeClass("layui-hide");
}