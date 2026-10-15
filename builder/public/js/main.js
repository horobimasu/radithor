import tag from './tag';
import site from './site';
import build from './build';

(async () => {
    await site.preSiteInit();
    await tag.initTag();
    await site.initSite();
    build.initBuild();
})();
