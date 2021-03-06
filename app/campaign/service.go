package campaign

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface{
	GetCampaigns(userID int)([]Campaign, error)
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	Update(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error)
	SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service{
	return &service{repository}
}

func (s *service) GetCampaigns(userID int)([]Campaign, error){
	if userID == 0 {
		campaigns, err := s.repository.FindAll()
		if err !=  nil {
			return campaigns, err
		}

		return campaigns, nil
	}
	
	campaigns, err := s.repository.FindByUserID(userID)
	if err !=  nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (s *service) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error){
	campaign, err := s.repository.FindByID(input.ID)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error){
	
	campaign := Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.GoalAmount = input.GoalAmount
	campaign.UserID = input.User.ID
	campaign.Perks = input.Perks
	// pembuatan slug
	slugCandidate := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(slugCandidate)

	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}

func (s *service) Update(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error){
	// cari campaignnya yang didapet dari uri
	campaign, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return campaign, err
	}

	// cek update pemilik dari campaignnya
	if campaign.UserID != inputData.User.ID {
		return campaign, errors.New("Not an owner of the campaign")
	}

	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.GoalAmount = inputData.GoalAmount
	campaign.Perks = inputData.Perks

	// update campaign
	updatedCampaign, err := s.repository.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}

	return updatedCampaign, nil
}

func (s *service) SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error){
	// ambil userID dari campaign
	campaign, err := s.repository.FindByID(input.CampaignID)
	if err != nil {
		return CampaignImage{}, err
	}

	// komparasikan user input dengan user campaignID
	if campaign.UserID != input.User.ID {
		return CampaignImage{}, errors.New("Not an owner of the campaign")
	}

	isPrimary := 0
	if input.IsPrimary == "true" {
		isPrimary = 1
		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignID)
		if err != nil {
			return CampaignImage{}, err
		}
	}

	// mapping
	campaignImage := CampaignImage{}
	campaignImage.CampaignID = input.CampaignID
	campaignImage.IsPrimary = isPrimary
	campaignImage.FileName = fileLocation

	newCampaign, err := s.repository.CreateImage(campaignImage)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil

}