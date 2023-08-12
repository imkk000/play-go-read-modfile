package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/mod/modfile"
)

func main() {
	content, err := RetrieveModFile("https://github.com/junegunn/fzf/blob/master/go.mod", GithubJSONParser)
	if err != nil {
		fmt.Println("retrieve mod file:", err)
		os.Exit(1)
	}
	f, err := modfile.Parse("go.mod", content, nil)
	if err != nil {
		fmt.Println("parse mod file:", err)
		os.Exit(1)
	}
	fmt.Println("go version:", f.Go.Version)
	fmt.Println("modules:")
	for _, r := range f.Require {
		fmt.Printf("- %s\n", ReplacePathVersion(r.Mod.Path))
	}
}

func ReplacePathVersion(path string) string {
	match, _ := regexp.MatchString(`(/v\d+)$`, path)
	if match {
		ss := strings.Split(path, "/")
		return strings.Join(ss[:len(ss)-1], "/")
	}
	return path
}

func RetrieveModFile(url string, customParser ...func(io.Reader) ([]byte, error)) ([]byte, error) {
	if customParser == nil {
		customParser = append(customParser, io.ReadAll)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "*/*")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	parser := customParser[0]
	return parser(resp.Body)
}

func GithubJSONParser(r io.Reader) ([]byte, error) {
	var body GithubResponseBody
	if err := json.NewDecoder(r).Decode(&body); err != nil {
		return nil, err
	}
	rawLines := body.Payload.Blob.RawLines
	if rawLines == nil {
		return nil, errors.New("empty content")
	}
	return []byte(strings.Join(rawLines, "\n")), nil
}

type GithubResponseBody struct {
	Payload struct {
		AllShortcutsEnabled bool `json:"allShortcutsEnabled,omitempty"`
		FileTree            struct {
			NAMING_FAILED struct {
				Items []struct {
					Name        string `json:"name,omitempty"`
					Path        string `json:"path,omitempty"`
					ContentType string `json:"contentType,omitempty"`
				} `json:"items,omitempty"`
				TotalCount int `json:"totalCount,omitempty"`
			} `json:",omitempty"`
		} `json:"fileTree,omitempty"`
		FileTreeProcessingTime float64 `json:"fileTreeProcessingTime,omitempty"`
		FoldersToFetch         []any   `json:"foldersToFetch,omitempty"`
		ReducedMotionEnabled   any     `json:"reducedMotionEnabled,omitempty"`
		Repo                   struct {
			ID                 int       `json:"id,omitempty"`
			DefaultBranch      string    `json:"defaultBranch,omitempty"`
			Name               string    `json:"name,omitempty"`
			OwnerLogin         string    `json:"ownerLogin,omitempty"`
			CurrentUserCanPush bool      `json:"currentUserCanPush,omitempty"`
			IsFork             bool      `json:"isFork,omitempty"`
			IsEmpty            bool      `json:"isEmpty,omitempty"`
			CreatedAt          time.Time `json:"createdAt,omitempty"`
			OwnerAvatar        string    `json:"ownerAvatar,omitempty"`
			Public             bool      `json:"public,omitempty"`
			Private            bool      `json:"private,omitempty"`
			IsOrgOwned         bool      `json:"isOrgOwned,omitempty"`
		} `json:"repo,omitempty"`
		SymbolsExpanded bool `json:"symbolsExpanded,omitempty"`
		TreeExpanded    bool `json:"treeExpanded,omitempty"`
		RefInfo         struct {
			Name         string `json:"name,omitempty"`
			ListCacheKey string `json:"listCacheKey,omitempty"`
			CanEdit      bool   `json:"canEdit,omitempty"`
			RefType      string `json:"refType,omitempty"`
			CurrentOid   string `json:"currentOid,omitempty"`
		} `json:"refInfo,omitempty"`
		Path        string `json:"path,omitempty"`
		CurrentUser any    `json:"currentUser,omitempty"`
		Blob        struct {
			RawLines          []string `json:"rawLines,omitempty"`
			StylingDirectives [][]struct {
				Start    int    `json:"start,omitempty"`
				End      int    `json:"end,omitempty"`
				CSSClass string `json:"cssClass,omitempty"`
			} `json:"stylingDirectives,omitempty"`
			Csv            any `json:"csv,omitempty"`
			CsvError       any `json:"csvError,omitempty"`
			DependabotInfo struct {
				ShowConfigurationBanner        bool   `json:"showConfigurationBanner,omitempty"`
				ConfigFilePath                 any    `json:"configFilePath,omitempty"`
				NetworkDependabotPath          string `json:"networkDependabotPath,omitempty"`
				DismissConfigurationNoticePath string `json:"dismissConfigurationNoticePath,omitempty"`
				ConfigurationNoticeDismissed   any    `json:"configurationNoticeDismissed,omitempty"`
				RepoAlertsPath                 string `json:"repoAlertsPath,omitempty"`
				RepoSecurityAndAnalysisPath    string `json:"repoSecurityAndAnalysisPath,omitempty"`
				RepoOwnerIsOrg                 bool   `json:"repoOwnerIsOrg,omitempty"`
				CurrentUserCanAdminRepo        bool   `json:"currentUserCanAdminRepo,omitempty"`
			} `json:"dependabotInfo,omitempty"`
			DisplayName string `json:"displayName,omitempty"`
			DisplayURL  string `json:"displayUrl,omitempty"`
			HeaderInfo  struct {
				BlobSize   string `json:"blobSize,omitempty"`
				DeleteInfo struct {
					DeleteTooltip string `json:"deleteTooltip,omitempty"`
				} `json:"deleteInfo,omitempty"`
				EditInfo struct {
					EditTooltip string `json:"editTooltip,omitempty"`
				} `json:"editInfo,omitempty"`
				GhDesktopPath    string `json:"ghDesktopPath,omitempty"`
				GitLfsPath       any    `json:"gitLfsPath,omitempty"`
				OnBranch         bool   `json:"onBranch,omitempty"`
				ShortPath        string `json:"shortPath,omitempty"`
				SiteNavLoginPath string `json:"siteNavLoginPath,omitempty"`
				IsCSV            bool   `json:"isCSV,omitempty"`
				IsRichtext       bool   `json:"isRichtext,omitempty"`
				Toc              any    `json:"toc,omitempty"`
				LineInfo         struct {
					TruncatedLoc  string `json:"truncatedLoc,omitempty"`
					TruncatedSloc string `json:"truncatedSloc,omitempty"`
				} `json:"lineInfo,omitempty"`
				Mode string `json:"mode,omitempty"`
			} `json:"headerInfo,omitempty"`
			Image                      bool   `json:"image,omitempty"`
			IsCodeownersFile           any    `json:"isCodeownersFile,omitempty"`
			IsPlain                    bool   `json:"isPlain,omitempty"`
			IsValidLegacyIssueTemplate bool   `json:"isValidLegacyIssueTemplate,omitempty"`
			IssueTemplateHelpURL       string `json:"issueTemplateHelpUrl,omitempty"`
			IssueTemplate              any    `json:"issueTemplate,omitempty"`
			DiscussionTemplate         any    `json:"discussionTemplate,omitempty"`
			Language                   string `json:"language,omitempty"`
			LanguageID                 int    `json:"languageID,omitempty"`
			Large                      bool   `json:"large,omitempty"`
			LoggedIn                   bool   `json:"loggedIn,omitempty"`
			NewDiscussionPath          string `json:"newDiscussionPath,omitempty"`
			NewIssuePath               string `json:"newIssuePath,omitempty"`
			PlanSupportInfo            struct {
				RepoIsFork                     any    `json:"repoIsFork,omitempty"`
				RepoOwnedByCurrentUser         any    `json:"repoOwnedByCurrentUser,omitempty"`
				RequestFullPath                string `json:"requestFullPath,omitempty"`
				ShowFreeOrgGatedFeatureMessage any    `json:"showFreeOrgGatedFeatureMessage,omitempty"`
				ShowPlanSupportBanner          any    `json:"showPlanSupportBanner,omitempty"`
				UpgradeDataAttributes          any    `json:"upgradeDataAttributes,omitempty"`
				UpgradePath                    any    `json:"upgradePath,omitempty"`
			} `json:"planSupportInfo,omitempty"`
			PublishBannersInfo struct {
				DismissActionNoticePath string `json:"dismissActionNoticePath,omitempty"`
				DismissStackNoticePath  string `json:"dismissStackNoticePath,omitempty"`
				ReleasePath             string `json:"releasePath,omitempty"`
				ShowPublishActionBanner bool   `json:"showPublishActionBanner,omitempty"`
				ShowPublishStackBanner  bool   `json:"showPublishStackBanner,omitempty"`
			} `json:"publishBannersInfo,omitempty"`
			RenderImageOrRaw bool `json:"renderImageOrRaw,omitempty"`
			RichText         any  `json:"richText,omitempty"`
			RenderedFileInfo any  `json:"renderedFileInfo,omitempty"`
			ShortPath        any  `json:"shortPath,omitempty"`
			TabSize          int  `json:"tabSize,omitempty"`
			TopBannersInfo   struct {
				OverridingGlobalFundingFile       bool   `json:"overridingGlobalFundingFile,omitempty"`
				GlobalPreferredFundingPath        any    `json:"globalPreferredFundingPath,omitempty"`
				RepoOwner                         string `json:"repoOwner,omitempty"`
				RepoName                          string `json:"repoName,omitempty"`
				ShowInvalidCitationWarning        bool   `json:"showInvalidCitationWarning,omitempty"`
				CitationHelpURL                   string `json:"citationHelpUrl,omitempty"`
				ShowDependabotConfigurationBanner bool   `json:"showDependabotConfigurationBanner,omitempty"`
				ActionsOnboardingTip              any    `json:"actionsOnboardingTip,omitempty"`
			} `json:"topBannersInfo,omitempty"`
			Truncated           bool `json:"truncated,omitempty"`
			Viewable            bool `json:"viewable,omitempty"`
			WorkflowRedirectURL any  `json:"workflowRedirectUrl,omitempty"`
			Symbols             struct {
				TimedOut    bool  `json:"timedOut,omitempty"`
				NotAnalyzed bool  `json:"notAnalyzed,omitempty"`
				Symbols     []any `json:"symbols,omitempty"`
			} `json:"symbols,omitempty"`
		} `json:"blob,omitempty"`
		CopilotInfo any `json:"copilotInfo,omitempty"`
		CsrfTokens  struct {
			JunegunnFzfBranches struct {
				Post string `json:"post,omitempty"`
			} `json:"/junegunn/fzf/branches,omitempty"`
			ReposPreferences struct {
				Post string `json:"post,omitempty"`
			} `json:"/repos/preferences,omitempty"`
		} `json:"csrf_tokens,omitempty"`
	} `json:"payload,omitempty"`
	Title string `json:"title,omitempty"`
}
