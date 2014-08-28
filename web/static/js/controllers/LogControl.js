function LogControl($scope, authService, resourcesService) {
    authService.checkLogin($scope);
    $scope.breadcrumbs = [
        { label: 'breadcrumb_logs', itemClass: 'active' }
    ];
    setInterval(function() {
        var logsframe = document.getElementById("logsframe");
        if (logsframe) {
            var h = logsframe.contentWindow.document.body.clientHeight;
            logsframe.height = h + "px";
        }
    }, 100);
}
