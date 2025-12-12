// This script patches Object.keys BEFORE any other code runs
// It must be loaded in index.html or app.vue before any imports
(function() {
  'use strict';
  if (typeof window === 'undefined') return;
  
  const originalObjectKeys = Object.keys;
  
  Object.keys = function(obj) {
    if (obj === null || obj === undefined) {
      return [];
    }
    if (typeof obj !== 'object' && typeof obj !== 'function') {
      return [];
    }
    try {
      return originalObjectKeys.call(this, obj);
    } catch (e) {
      return [];
    }
  };
  
  console.log('[Object.keys Patch] Applied successfully');
})();

