// Navigation utility - setup navigation based on user role
function setupNavigation() {
    const user = getUser();
    const navMenus = document.querySelectorAll('[data-nav-menu]');
    
    // Get current page to highlight active link
    const currentPath = window.location.pathname;
    
    navMenus.forEach(navMenu => {
        if (user && user.role === 'caregiver') {
            // Caregiver navigation
            const isDashboard = currentPath === '/dashboard.html' || currentPath === '/';
            const isProfile = currentPath === '/caregiver-profile-setup.html';
            const isAvailability = currentPath === '/caregiver-availability.html';
            const isLicenses = currentPath === '/caregiver-licenses.html';
            const isContracts = currentPath === '/contracts.html';
            
            const isUserProfile = currentPath === '/profile.html';
            navMenu.innerHTML = `
                <a href="/dashboard.html" class="${isDashboard ? 'text-primary-600 border-b-2 border-primary-600' : 'text-gray-700 hover:text-primary-600'} px-3 py-2 rounded-md text-sm font-medium">儀表板</a>
                <a href="/caregiver-profile-setup.html" class="${isProfile ? 'text-primary-600 border-b-2 border-primary-600' : 'text-gray-700 hover:text-primary-600'} px-3 py-2 rounded-md text-sm font-medium">我的檔案</a>
                <a href="/caregiver-availability.html" class="${isAvailability ? 'text-primary-600 border-b-2 border-primary-600' : 'text-gray-700 hover:text-primary-600'} px-3 py-2 rounded-md text-sm font-medium">服務時間</a>
                <a href="/caregiver-licenses.html" class="${isLicenses ? 'text-primary-600 border-b-2 border-primary-600' : 'text-gray-700 hover:text-primary-600'} px-3 py-2 rounded-md text-sm font-medium">證照管理</a>
                <a href="/contracts.html" class="${isContracts ? 'text-primary-600 border-b-2 border-primary-600' : 'text-gray-700 hover:text-primary-600'} px-3 py-2 rounded-md text-sm font-medium">我的合約</a>
                <a href="/profile.html" class="${isUserProfile ? 'text-primary-600 border-b-2 border-primary-600' : 'text-gray-700 hover:text-primary-600'} px-3 py-2 rounded-md text-sm font-medium">個人資訊</a>
            `;
        } else {
            // User navigation
            const isDashboard = currentPath === '/dashboard.html' || currentPath === '/';
            const isSearch = currentPath === '/caregiver-search.html';
            const isContracts = currentPath === '/contracts.html';
            
            const isUserProfile = currentPath === '/profile.html';
            navMenu.innerHTML = `
                <a href="/dashboard.html" class="${isDashboard ? 'text-primary-600 border-b-2 border-primary-600' : 'text-gray-700 hover:text-primary-600'} px-3 py-2 rounded-md text-sm font-medium">儀表板</a>
                <a href="/caregiver-search.html" class="${isSearch ? 'text-primary-600 border-b-2 border-primary-600' : 'text-gray-700 hover:text-primary-600'} px-3 py-2 rounded-md text-sm font-medium">搜尋照護者</a>
                <a href="/contracts.html" class="${isContracts ? 'text-primary-600 border-b-2 border-primary-600' : 'text-gray-700 hover:text-primary-600'} px-3 py-2 rounded-md text-sm font-medium">我的合約</a>
                <a href="/profile.html" class="${isUserProfile ? 'text-primary-600 border-b-2 border-primary-600' : 'text-gray-700 hover:text-primary-600'} px-3 py-2 rounded-md text-sm font-medium">個人資訊</a>
            `;
        }
    });
}

// Call setupNavigation when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', setupNavigation);
} else {
    setupNavigation();
}

